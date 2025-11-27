package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)



func HashPassword(password string)(string, error){
	hashed, err := argon2id.CreateHash(password,argon2id.DefaultParams)
	if err != nil {
		return "", err
	}	
	return hashed, nil
}

func VerifyPassword(password, hashed string)(bool, error){
	match, err := argon2id.ComparePasswordAndHash(password,hashed)
	if err != nil {
		return false,err
	}
	return match, nil
}


func GetBearerToken(h http.Header)(string,error){
	auth := h.Get("Authorization")
	if auth == ""{
		return auth, fmt.Errorf("no token was provided")
	}
	token := strings.TrimPrefix(auth,"Bearer ")
	
	return token, nil
}


func GenerateJwtToken (secret string,id uuid.UUID, expiresIn time.Duration)(string, error){
	claims := jwt.RegisteredClaims{
		Issuer: "godocs",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Subject: id.String(),	
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,&claims)
	jwt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func ValidateJwt(secretString, tokenString string)(uuid.UUID, error){
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token)(any, error){
			return []byte(secretString), nil
		})
		if err != nil{
			return uuid.Nil, err
		}
		idString, err := token.Claims.GetSubject()
		if err != nil {
			return uuid.Nil, err
		}
		id, err := uuid.Parse(idString)
		if err != nil {
			return uuid.Nil, err
		}
		return id, nil
}

func GenerateRefreshToken() ( string, error ) {
	b := make([]byte,32)
	rand.Read(b)
	token := hex.EncodeToString(b)

	return string(token), nil


}