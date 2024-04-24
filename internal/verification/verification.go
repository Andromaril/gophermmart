package verification

import (
	"fmt"
	"time"

	"github.com/andromaril/gophermmart/internal/errormart"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

const TokenExp = time.Hour * 3
const SecretKey = "supersecretkey"

func main() {
	tokenString, err := BuildJWTString()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tokenString)
}

func BuildJWTString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: 1,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		e := errormart.NewMartError(err)
		return "", fmt.Errorf("error %q", e.Error())
	}

	return tokenString, nil
}

func GetUserID(tokenString string) int {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})
	if err != nil {
		e := errormart.NewMartError(err)
		log.Error(e.Error())
		return -1
	}

	if !token.Valid {
		return -1
	}
	return claims.UserID
}
