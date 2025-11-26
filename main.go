package main

import (
	"database/sql"
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
	if dbUrl == ""{
		log.Fatal("Missing database Url")
	}
	if secretKey == ""{
		log.Fatal("Missing secret Key")
	}
	db, err := sql.Open("postgres",dbUrl)

	if err != nil {
		log.Fatal("Error while opening database connexion")
	}
	
	dbQueries := database.New(db)

	cfg := config.Config{
		Port: port,
	}
	apiCfg := api.ApiConfig{
		Db:dbQueries,
		SecretKey: secretKey,
	}


	mux := http.NewServeMux()
	srv := &http.Server{
		Addr: ":" + cfg.Port ,
		Handler: mux,
	}

	mux.HandleFunc("POST /api/users",apiCfg.CreateUser)
	mux.HandleFunc("POST /api/auth/login",apiCfg.LoginUser)



	fmt.Printf("Serving on:  http://localhost:%s\n", cfg.Port)
	log.Fatal(srv.ListenAndServe())

}