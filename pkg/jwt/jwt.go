package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"murakali/config"
	"net/http"
	"strings"
	"time"
)

type AccessClaims struct {
	ID     string `json:"id"`
	RoleID int    `json:"role_id"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

type RegisterClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWTAccessToken(userID string, userRole int, config *config.Config) (string, error) {
	claims := &AccessClaims{
		ID:     userID,
		RoleID: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.JWT.AccessExpMin) * time.Minute)),
			Issuer:    config.JWT.JwtIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.JWT.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateJWTRefreshToken(userID string, config *config.Config) (string, error) {
	claims := &RefreshClaims{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.JWT.RefreshExpMin) * time.Minute)),
			Issuer:    config.JWT.JwtIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.JWT.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateJWTRegisterToken(email string, config *config.Config) (string, error) {
	claims := &RegisterClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.JWT.RefreshExpMin) * time.Minute)),
			Issuer:    config.JWT.JwtIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.JWT.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractJWT(tokenString, jwtKey string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token ")
	}

	return claims, nil
}

func ExtractJWTFromRequest(r *http.Request, jwtKey string) (map[string]interface{}, error) {
	tokenString := ExtractBearerToken(r)

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token ")
	}

	return claims, nil
}

func ExtractBearerToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	token := strings.Split(bearerToken, " ")
	if len(token) == 2 {
		return token[1]
	}

	return ""
}
