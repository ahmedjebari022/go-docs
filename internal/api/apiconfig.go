package api

import "github.com/ahmedjebari022/go-docs/internal/database"


type ApiConfig struct{
	Db 		*database.Queries
	SecretKey string
	CookieKey []byte
	AssetsPath string
	Port string
}