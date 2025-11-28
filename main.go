package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ahmedjebari022/go-docs/internal/api"
	"github.com/ahmedjebari022/go-docs/internal/config"
	"github.com/ahmedjebari022/go-docs/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main(){
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	secretKey := os.Getenv("SECRET_KEY")
	port := os.Getenv("PORT")
	cookieKey := os.Getenv("COOKIE_SECRET")
	assetsPath := os.Getenv("ASSETS_ROOT")
	if dbUrl == ""{
		log.Fatal("Missing database Url")
	}
	if secretKey == ""{
		log.Fatal("Missing secret Key")
	}
	if cookieKey == ""{
		log.Fatal("Missing cookie secret Key")
	}
	if assetsPath == ""{
		log.Fatal("Missing assets path")
	}
	db, err := sql.Open("postgres",dbUrl)

	if err != nil {
		log.Fatal("Error while opening database connexion")
	}
	
	dbQueries := database.New(db)
	ck, err := hex.DecodeString(cookieKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	cfg := config.Config{
		Port: port,
	}
	apiCfg := api.ApiConfig{
		DbC: db,
		Db:dbQueries,
		SecretKey: secretKey,
		CookieKey: ck,
		AssetsPath: assetsPath,
		Port: port,
	}


	mux := http.NewServeMux()
	srv := &http.Server{
		Addr: ":" + cfg.Port ,
		Handler: mux,
	}

	mux.HandleFunc("POST /api/users",apiCfg.CreateUser)
	mux.HandleFunc("POST /api/auth/login",apiCfg.LoginUser)
	mux.HandleFunc("GET /api/cookie",apiCfg.ReaderCookieHandler)
	mux.HandleFunc("POST /api/cookie/refresh",apiCfg.RefreshTokenHandler)
	mux.Handle("POST /api/documents",apiCfg.AuthMiddleware(http.HandlerFunc(apiCfg.CreateDocumentHandler)))



	fmt.Printf("Serving on:  http://localhost:%s\n", cfg.Port)
	log.Fatal(srv.ListenAndServe())

}