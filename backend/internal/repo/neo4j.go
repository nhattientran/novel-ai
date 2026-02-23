package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Driver wraps Neo4j driver with helper methods
type Driver struct {
	driver neo4j.DriverWithContext
}

// NewDriver creates a new Neo4j driver instance
func NewDriver(uri, user, password string) (*Driver, error) {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(user, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	return &Driver{driver: driver}, nil
}

// Close gracefully closes the driver
func (d *Driver) Close(ctx context.Context) error {
	return d.driver.Close(ctx)
}

// VerifyConnectivity checks if the driver can connect to Neo4j
func (d *Driver) VerifyConnectivity(ctx context.Context) error {
	return d.driver.VerifyConnectivity(ctx)
}

// ExecuteRead executes a read transaction with timeout and returns records
func (d *Driver) ExecuteRead(ctx context.Context, query string, params map[string]any) ([]*neo4j.Record, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	session := d.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return result.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}

	return result.([]*neo4j.Record), nil
}

// ExecuteWrite executes a write transaction with timeout and returns records
func (d *Driver) ExecuteWrite(ctx context.Context, query string, params map[string]any) ([]*neo4j.Record, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	session := d.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		return result.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}

	return result.([]*neo4j.Record), nil
}

// RunQuery executes a simple query without transaction management
func (d *Driver) RunQuery(ctx context.Context, query string, params map[string]any) (neo4j.ResultWithContext, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	session := d.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	return session.Run(ctx, query, params)
}
