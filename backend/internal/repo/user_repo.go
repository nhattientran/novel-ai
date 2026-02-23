package repo

import (
	"context"
	"fmt"

	"novel-ai/internal/domain"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// UserRepo handles user data access
type UserRepo struct {
	driver *Driver
}

// NewUserRepo creates a new user repository
func NewUserRepo(driver *Driver) *UserRepo {
	return &UserRepo{driver: driver}
}

// Create creates a new user in Neo4j
func (r *UserRepo) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		CREATE (u:User {
			id: $id,
			username: $username,
			email: $email,
			password_hash: $password_hash,
			role: $role
		})
		RETURN u { .id, .username, .email, .role } AS user
	`

	params := map[string]any{
		"id":            user.ID,
		"username":      user.Username,
		"email":         user.Email,
		"password_hash": user.PasswordHash,
		"role":          user.Role,
	}

	records, err := r.driver.ExecuteWrite(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no user created")
	}

	return r.recordToUser(records[0], false)
}

// FindByEmail finds a user by email (includes password hash for auth)
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		MATCH (u:User {email: $email})
		RETURN u { .id, .username, .email, .password_hash, .role } AS user
	`

	params := map[string]any{
		"email": email,
	}

	records, err := r.driver.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	if len(records) == 0 {
		return nil, nil // user not found
	}

	return r.recordToUser(records[0], true)
}

// FindByID finds a user by ID (excludes password hash)
func (r *UserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		MATCH (u:User {id: $id})
		RETURN u { .id, .username, .email, .role } AS user
	`

	params := map[string]any{
		"id": id,
	}

	records, err := r.driver.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	if len(records) == 0 {
		return nil, nil // user not found
	}

	return r.recordToUser(records[0], false)
}

// recordToUser converts a Neo4j record to User domain model
func (r *UserRepo) recordToUser(record *neo4j.Record, includePassword bool) (*domain.User, error) {
	userMap, ok := record.Get("user")
	if !ok {
		return nil, fmt.Errorf("user not found in record")
	}

	userData, ok := userMap.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid user data type")
	}

	user := &domain.User{
		ID:       getString(userData, "id"),
		Username: getString(userData, "username"),
		Email:    getString(userData, "email"),
		Role:     getString(userData, "role"),
	}

	if includePassword {
		user.PasswordHash = getString(userData, "password_hash")
	}

	return user, nil
}

// getString safely extracts a string value from a map
func getString(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
