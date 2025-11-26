package api

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/ahmedjebari022/go-docs/internal/auth"
	"github.com/ahmedjebari022/go-docs/internal/database"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)



func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request){
		defer r.Body.Close()
		type reqBody struct{
			Email string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required,password"`
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
		
		validate := validator.New(validator.WithRequiredStructEnabled())
		validate.RegisterValidation("password", func(fl validator.FieldLevel)bool{
			password := fl.Field().String()
			hasMinLength := len(password) >= 8
			hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
			hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
			hasDigit := regexp.MustCompile(`\d`).MatchString(password)
			hasSpecial := regexp.MustCompile(`[@$!%*?&]`).MatchString(password)
			return hasUpper && hasLower && hasDigit && hasSpecial && hasMinLength
		})
		err = validate.Struct(params)
		if err != nil {
			
			// var validateErrors validator.ValidationErrors
			// errorsMsg := ""
			// if errors.As(err, &validateErrors){
			// 	for _, e := range validateErrors{
			// 		errorsMsg += fmt.Sprintf("%s, %s",e.Field(),e.Error())
			// 		if e.Tag() == "password"{
			// 			errorsMsg += "password : must contain at least a number an uppercase and a special chars"
			// 		}
			// 	}
			// }
			respondWithError(w,400,err.Error())
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

func (cfg *ApiConfig) LoginUser(w http.ResponseWriter, r *http.Request){
		type RequestBody struct{
			Email  string 	`json:"email"` 
			Password string `json:"password"`
		}
		type ResponseBody struct{
			Email string `json:"email"`
			Password string `json:"password"`
			Token 		string `json:"token"`
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithError(w,http.StatusBadRequest,err.Error())
			return
		}
		var params RequestBody
		err = json.Unmarshal(body,&params)
		if err != nil {
			respondWithError(w,500,err.Error())
			return
		}
		user, err := cfg.Db.GetUserByEmail(r.Context(),params.Email)
		if err != nil {
			respondWithError(w,400,err.Error())
			return
		}
		if match, _ := auth.VerifyPassword(params.Password,user.HashedPassword); !match{
			respondWithError(w,http.StatusBadRequest,"invalid crredientials")
			return
		}
		jwt, err := auth.GenerateJwtToken(cfg.SecretKey,user.ID,time.Hour)
		if err != nil {
			respondWithError(w,500,err.Error())
		}

		res := ResponseBody{
			Email: user.Email,
			Password: params.Password,
			Token: jwt,
		}
		RespondWithJson(w,http.StatusOK,res)
}

