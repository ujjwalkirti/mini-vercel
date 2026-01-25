package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"sync"
	"time"
)

type JWKSCache struct {
	mu        sync.RWMutex
	keys      map[string]interface{}
	expiresAt time.Time
	jwksURL   string
	ttl       time.Duration
}

func NewJWKSCache(jwksURL string, ttl time.Duration) *JWKSCache {
	return &JWKSCache{
		keys:    make(map[string]interface{}),
		jwksURL: jwksURL,
		ttl:     ttl,
	}
}

func (c *JWKSCache) Get(kid string) (interface{}, error) {
	c.mu.RLock()
	if time.Now().Before(c.expiresAt) {
		if key, ok := c.keys[kid]; ok {
			c.mu.RUnlock()
			return key, nil
		}
	}
	c.mu.RUnlock()

	// Cache expired or key missing → refresh
	if err := c.refresh(); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	key, ok := c.keys[kid]
	if !ok {
		return nil, errors.New("public key not found after refresh")
	}

	return key, nil
}

func (c *JWKSCache) refresh() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check (another goroutine may have refreshed)
	if time.Now().Before(c.expiresAt) {
		return nil
	}

	resp, err := http.Get(c.jwksURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
	}

	keys := make(map[string]interface{})

	for _, k := range jwks.Keys {
		switch k.Kty {

		case "RSA":
			nBytes, err := base64.RawURLEncoding.DecodeString(k.N)
			if err != nil {
				continue
			}

			eBytes, err := base64.RawURLEncoding.DecodeString(k.E)
			if err != nil {
				continue
			}

			n := new(big.Int).SetBytes(nBytes)

			e := 0
			for _, b := range eBytes {
				e = e<<8 + int(b)
			}

			keys[k.Kid] = &rsa.PublicKey{
				N: n,
				E: e,
			}

		case "EC":
			if k.Crv != "P-256" {
				continue
			}

			xBytes, err := base64.RawURLEncoding.DecodeString(k.X)
			if err != nil {
				continue
			}

			yBytes, err := base64.RawURLEncoding.DecodeString(k.Y)
			if err != nil {
				continue
			}

			pub := &ecdsa.PublicKey{
				Curve: elliptic.P256(),
				X:     new(big.Int).SetBytes(xBytes),
				Y:     new(big.Int).SetBytes(yBytes),
			}

			keys[k.Kid] = pub
		}
	}

	c.keys = keys
	c.expiresAt = time.Now().Add(c.ttl)
	return nil
}
