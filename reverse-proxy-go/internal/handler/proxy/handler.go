package proxy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"reverse-proxy/internal/repository/project"
)

type Handler struct {
	repo        *project.Repository
	r2PublicURL string
}

func NewHandler(repo *project.Repository, r2PublicURL string) *Handler {
	return &Handler{
		repo:        repo,
		r2PublicURL: r2PublicURL,
	}
}

// ProxyRequest handles all incoming requests and proxies them to R2
func (h *Handler) ProxyRequest(w http.ResponseWriter, r *http.Request) {
	// Extract subdomain from hostname
	hostname := r.Host
	subdomain := h.extractSubdomain(hostname)
	log.Printf("Received request for subdomain %s", subdomain)
	if subdomain == "" {
		http.Error(w, "Invalid subdomain", http.StatusBadRequest)
		return
	}

	// Query database for project with subdomain
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	proj, err := h.repo.FindBySubdomain(ctx, subdomain)
	if err != nil {
		log.Printf("Project not found for subdomain %s: %v", subdomain, err)
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Check if there's a ready deployment
	if len(proj.Deployments) == 0 {
		log.Printf("No deployments found for project %s", proj.ID)
		http.Error(w, "No deployment available", http.StatusNotFound)
		return
	}

	// Build the target URL for R2
	targetPath := r.URL.Path
	if targetPath == "/" || targetPath == "" {
		targetPath = "/index.html"
	}

	targetURL := fmt.Sprintf("%s/%s%s", h.r2PublicURL, proj.ID, targetPath)
	log.Printf("Proxying request: %s -> %s", r.URL.Path, targetURL)

	// Parse the target URL
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Printf("Failed to parse target URL: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Modify the request
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
	}

	// Error handler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		http.Error(w, "Bad gateway", http.StatusBadGateway)
	}

	// Serve the proxied request
	proxy.ServeHTTP(w, r)
}

// extractSubdomain extracts the subdomain from the hostname
// Example: myapp.localhost:8001 -> myapp
func (h *Handler) extractSubdomain(hostname string) string {
	// Remove port if present
	host := hostname
	if idx := strings.Index(hostname, ":"); idx != -1 {
		host = hostname[:idx]
	}

	// Split by dots
	parts := strings.Split(host, ".")
	if len(parts) > 1 {
		// Return the first part as subdomain
		return parts[0]
	}

	// If no subdomain, return empty
	return ""
}
