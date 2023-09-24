package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var TokenName = "GORRC_Token"
var tempSecretKey = []byte("TempSecretKey") //TODO: Load from env variable
var AuthTime = time.Now().Add(24 * time.Hour)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(userName string) (string, error) {
	claims := &Claims{
		Username: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(AuthTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(tempSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed using secret key: %w", err)
	}
	return tokenString, nil
}

func ValidateJWT(token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return tempSecretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("failed validating jwt token: %w", err)
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims["username"].(string), nil
	}
	return "", fmt.Errorf("failed token validation")
}

func ValidateCookie(req *http.Request) (string, error) {
	cookie, err := req.Cookie(TokenName)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		return "", fmt.Errorf("failed parsing cookie for token: %w", err)
	}
	if cookie.Valid() != nil {
		return "", fmt.Errorf("no cookie found with token")
	}

	userName, err := ValidateJWT(cookie.Value)
	if err != nil {
		return "", fmt.Errorf("cookie failed token validation: %w", err)
	}
	return userName, nil
}
