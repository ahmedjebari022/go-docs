package auth

import "github.com/alexedwards/argon2id"



func HashPassword(password string)(string, error){
	hashed, err := argon2id.CreateHash(password,argon2id.DefaultParams)
	if err != nil {
		return "", err
	}	
	return hashed, nil
}

func VerifyPassword(password, hashed string)(bool, error){
	match, err := argon2id.ComparePasswordAndHash(password,hashed)
	if err == nil {
		return false,err
	}
	return match, nil
}