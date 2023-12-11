package token

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type MyJWTClaims struct {
	*jwt.RegisteredClaims
	UserInfo interface{}
}

type UserJWTData struct {
	Username string
	Groups   []string
	Admin    bool
}

var privateKey, publicKey, keyReadErr = ReadKeyFiles()

func Create(sub string, userInfo interface{}) (string, error) {
	if keyReadErr != nil {
		return "", fmt.Errorf("create: read key files: %w", keyReadErr)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	exp := time.Now().Add(time.Hour * 12)

	claims := &MyJWTClaims{
		&jwt.RegisteredClaims{
			Subject:   sub,
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		userInfo,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func GetClaimsFromToken(tokenString string) (jwt.MapClaims, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, fmt.Errorf("get claims: parse key: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func ReadKeyFiles() ([]byte, []byte, error) {
	prvKey, err := os.ReadFile(os.Getenv("JWT_PRIVATE_KEY"))
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	pubKey, err := os.ReadFile(os.Getenv("JWT_PUBLIC_KEY"))
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	return prvKey, pubKey, nil
}
