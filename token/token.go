package token

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type MyJWTClaims struct {
	*jwt.RegisteredClaims
	UserInfo interface{}
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

	token := jwt.New(jwt.SigningMethodRS256)
	exp := time.Now().Add(time.Hour * 12)

	token.Claims = &MyJWTClaims{
		&jwt.RegisteredClaims{
			Subject:   sub,
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		userInfo,
	}

	val, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return val, nil
}

func GetClaimsFromToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
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

func SetJWTClaimsContext(c *gin.Context, claims jwt.MapClaims) {
	c.Set("claims", claims)
}

func JWTClaimsFromContext(c *gin.Context) jwt.MapClaims {
	return c.MustGet("claims").(jwt.MapClaims)
}
