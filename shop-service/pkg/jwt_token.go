package pkg

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

const (
	AccessTokenDuration  = 1 * time.Hour
	RefreshTokenDuration = 7 * 24 * time.Hour
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // in seconds
}

type TokenClaims struct {
	jwt.RegisteredClaims
	// Standard JWT Claims
	Sub   string   `json:"sub"`           // Subject (usually user ID)
	Iss   string   `json:"iss"`           // Issuer
	Aud   []string `json:"aud"`           // Audience
	Jti   string   `json:"jti,omitempty"` // JWT ID
	Scope string   `json:"scope"`         // OAuth 2.0 scope
}

func GenerateJWTToken(claims TokenClaims, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ValidateToken(tokenString, secretKey string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func GenerateTokenPair(userId, issuer string, audience []string, scope string, secretKey string) (*TokenPair, error) {
	// Generate Access Token
	accessClaims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        "at-" + uuid.New().String(),
		},
		Sub:   userId,
		Iss:   issuer,
		Aud:   audience,
		Scope: scope,
	}

	accessToken, err := GenerateJWTToken(accessClaims, secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate Refresh Token
	refreshClaims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        "rt-" + uuid.New().String(),
		},
		Sub:   userId,
		Iss:   issuer,
		Aud:   audience,
		Scope: "refresh_token",
	}

	refreshToken, err := GenerateJWTToken(refreshClaims, secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(AccessTokenDuration.Seconds()),
	}, nil
}

// Function to refresh tokens using a valid refresh token
func RefreshTokenPair(refreshToken string, secretKey string) (*TokenPair, error) {
	// Validate refresh token
	claims, err := ValidateToken(refreshToken, secretKey)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Verify it's a refresh token
	if claims.Scope != "refresh_token" {
		return nil, fmt.Errorf("invalid token type")
	}

	// Generate new token pair
	return GenerateTokenPair(
		claims.Sub,
		claims.Iss,
		claims.Aud,
		claims.Scope,
		secretKey,
	)
}
