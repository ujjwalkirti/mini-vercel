package client

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/ujjwalkirti/mini-vercel-api-server/internal/config"
)

// ClickHouseClient wraps the ClickHouse database connection
type ClickHouseClient struct {
	conn driver.Conn
}

// NewClickHouseClient creates a new ClickHouse client with configuration from config package
func NewClickHouseClient() (*ClickHouseClient, error) {
	conn, err := connect()
	if err != nil {
		return nil, fmt.Errorf("failed to create clickhouse client: %w", err)
	}

	return &ClickHouseClient{
		conn: conn,
	}, nil
}

// GetConn returns the underlying ClickHouse connection
func (c *ClickHouseClient) GetConn() driver.Conn {
	return c.conn
}

// Close closes the ClickHouse connection
func (c *ClickHouseClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Ping checks the connection to ClickHouse
func (c *ClickHouseClient) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx)
}

func connect() (driver.Conn, error) {
	ctx := context.Background()

	cfg := config.GetClickHouseConfig()
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "mini-vercel-api", Version: "1.0"},
			},
		},
		Debugf: func(format string, v ...interface{}) {
			fmt.Printf(format, v...)
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			return nil, fmt.Errorf("clickhouse exception [%d] %s: %s", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return conn, nil
}
