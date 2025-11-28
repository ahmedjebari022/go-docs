package api

import (
	"database/sql"

	"github.com/ahmedjebari022/go-docs/internal/database"
)


type ApiConfig struct{
	DbC 	*sql.DB
	Db 		*database.Queries
	SecretKey string
	CookieKey []byte
	AssetsPath string
	Port string
}