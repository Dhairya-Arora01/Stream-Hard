package auth

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Type for Public Claim along with Registered Claims.
type customClaims struct {
	UId int
	jwt.StandardClaims
}

// GenerateToken creates and return the token as string.
func GenerateToken(uid int) string {
	claims := customClaims{
		UId: uid,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "StreamHard",
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("saini72"))

	if err != nil {
		log.Println("Error while signing the token:", err)
	}

	return tokenString
}

// ValidateToken validates the user sent token
func ValidateToken(tokenString string) int {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("saini72"), nil
	})

	if err != nil || !token.Valid {
		log.Println("Invalid token")
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok {
		log.Println("Failed to get claims")
		return 0
	} else {
		return claims.UId
	}
}
