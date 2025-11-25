package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ahmedjebari022/go-docs/internal/auth"
	"github.com/ahmedjebari022/go-docs/internal/database"
	"github.com/google/uuid"
)


func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request){
		defer r.Body.Close()
		type reqBody struct{
			Email string `json:"email"`
			Password string `json:"password"`
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithError(w,500,err.Error())
			return
		}

		var params reqBody
		err = json.Unmarshal(data,&params)
		if err != nil {
			respondWithError(w,500,err.Error())
			return
		}
		
		if params.Email == "" || params.Password == ""{
			respondWithError(w,400,fmt.Sprint("missing required fields email or password"))
			return
		}

		hashed, err := auth.HashPassword(params.Password)
		if err != nil {
			respondWithError(w,500,err.Error())
			return 
		}

		user, err := cfg.Db.CreateUser(r.Context(), database.CreateUserParams{
			ID: uuid.New(),
			HashedPassword: hashed,
			Email: params.Email,
		})
		if err != nil {
			respondWithError(w,500,err.Error())
			return 
		}

		type ResponseBody struct{
			Email 			string `json:"email"`
			CreateAt 		time.Time `json:"created_at"` 
			UpdatedAt 		time.Time `json:"updated_at"`
			Password 		string `json:"password"`
		}

		RespondWithJson(w,200,ResponseBody{
			Email: user.Email,
			UpdatedAt: user.UpdatedAt,
			CreateAt: user.CreatedAt,
			Password: params.Password,
		})





}