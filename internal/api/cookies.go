package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
)
 
var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("cookie value is invalid")
)


func Write(w http.ResponseWriter, cookie http.Cookie) error {

	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))
	if len(cookie.String()) > 4096{
		return ErrValueTooLong

	}
	http.SetCookie(w,&cookie)
	return nil
}



func Read(r *http.Request, name string) (string, error) {
	c, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	value, err := base64.URLEncoding.DecodeString(c.Value)
	if err != nil {
		return "", ErrInvalidValue
	}
	return string(value), nil 
}




func WriteSigned(w http.ResponseWriter, cookie http.Cookie, key []byte) error {
	mac := hmac.New(sha256.New,key)
	mac.Write([]byte(cookie.Value))
	mac.Write([]byte(cookie.Name))
	signature := mac.Sum(nil)
	cookie.Value = string(signature) + cookie.Value
	
	return Write(w, cookie)
}

func ReadSigned(r *http.Request, name string, key []byte) (string, error){
	signedValue, err := Read(r, name)
	if err != nil {
		return "", err
	}
	if len(signedValue) < sha256.Size{
		return "", ErrInvalidValue
	}
	signature := signedValue[:sha256.Size]
	value := signedValue[sha256.Size:]
	
	mac := hmac.New(sha256.New,key)
	mac.Write([]byte(value))
	mac.Write([]byte(name))	
	signature2 := mac.Sum(nil)
	if match := hmac.Equal([]byte(signature),signature2) ; !match{
		return "", ErrInvalidValue
	}
	return value, nil
}
