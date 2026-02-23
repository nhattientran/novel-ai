package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	cookieName     = "access_token"
	cookieMaxAge   = 7 * 24 * 60 * 60 // 7 days in seconds
	contextUserID  = "userID"
	contextUserRole = "userRole"
)

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// CheckPassword compares a password with a hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken creates a new JWT token
func GenerateToken(userID, role, secret string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns claims
func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

// SetAuthCookie sets the JWT token in an HttpOnly cookie
func SetAuthCookie(c *gin.Context, token string) {
	c.SetCookie(
		cookieName,
		token,
		cookieMaxAge,
		"/",
		"",
		false, // Secure - set to true in production with HTTPS
		true,  // HttpOnly
	)
}

// ClearAuthCookie clears the auth cookie
func ClearAuthCookie(c *gin.Context) {
	c.SetCookie(
		cookieName,
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
}

// ExtractTokenFromCookie extracts JWT from cookie
func ExtractTokenFromCookie(c *gin.Context) (string, error) {
	token, err := c.Cookie(cookieName)
	if err != nil {
		return "", errors.New("auth cookie not found")
	}
	return token, nil
}

// SetUserContext sets user info in gin context
func SetUserContext(c *gin.Context, userID, role string) {
	c.Set(contextUserID, userID)
	c.Set(contextUserRole, role)
}

// GetUserID gets user ID from context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(contextUserID)
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}

// GetUserRole gets user role from context
func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get(contextUserRole)
	if !exists {
		return "", false
	}
	r, ok := role.(string)
	return r, ok
}

// RequireAuth middleware checks if user is authenticated
func RequireAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := ExtractTokenFromCookie(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		claims, err := ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		SetUserContext(c, claims.UserID, claims.Role)
		c.Next()
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := GetUserRole(c)
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		if userRole != role {
			c.JSON(403, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}
