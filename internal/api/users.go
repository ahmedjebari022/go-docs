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

const (
	refreshCookieName = "refreshCookie"
	accessCookieName  = "accessCookie"
)

func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type reqBody struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,password"`
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	var params reqBody
	err = json.Unmarshal(data, &params)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
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
		RespondWithError(w, 400, err.Error())
		return
	}
	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	user, err := cfg.Db.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		HashedPassword: hashed,
		Email:          params.Email,
	})
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	type ResponseBody struct {
		Email     string    `json:"email"`
		CreateAt  time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	RespondWithJson(w, 200, ResponseBody{
		Email:     user.Email,
		UpdatedAt: user.UpdatedAt,
		CreateAt:  user.CreatedAt,
	})
}

func (cfg *ApiConfig) LoginUser(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type ResponseBody struct {
		Email        string `json:"email"`
		RefreshToken string `json:"refresh_token"`
		Accestoken   string `json:"acces_token"`
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	var params RequestBody
	err = json.Unmarshal(body, &params)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	user, err := cfg.Db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}
	if match, _ := auth.VerifyPassword(params.Password, user.HashedPassword); !match {
		RespondWithError(w, http.StatusBadRequest, "invalid crredientials")
		return
	}
	jwt, err := auth.GenerateJwtToken(cfg.SecretKey, user.ID, time.Minute*15)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	accesCookie := http.Cookie{
		Name:     accessCookieName,
		Path:     "/",
		Domain:   "",
		Value:    jwt,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	refreshToken, err := auth.GenerateRefreshToken()

	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	rt, err := cfg.Db.CreateToken(r.Context(), database.CreateTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	refreshCookie := http.Cookie{
		Name:     refreshCookieName,
		Path:     "/",
		Domain:   "",
		Value:    rt.Token,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  rt.ExpiresAt,
		Secure:   true,
	}

	err = WriteSigned(w, refreshCookie, cfg.CookieKey)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	err = WriteSigned(w, accesCookie, cfg.CookieKey)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	res := ResponseBody{
		Email:        user.Email,
		Accestoken:   jwt,
		RefreshToken: refreshToken,
	}
	RespondWithJson(w, http.StatusOK, res)
}

func (cfg *ApiConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	value, err := ReadSigned(r, refreshCookieName, cfg.CookieKey)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	user, err := cfg.Db.GetUserByRefreshToken(r.Context(), value)
	if err != nil {
		RespondWithError(w, http.StatusForbidden, err.Error())
		return
	}
	jwt, err := auth.GenerateJwtToken(cfg.SecretKey, user.ID, time.Minute*15)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	cookie := http.Cookie{
		Name:     accessCookieName,
		Value:    jwt,
		Expires:  time.Now().Add(time.Minute * 15),
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}
	err = WriteSigned(w, cookie, cfg.CookieKey)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
}

func (cfg *ApiConfig) RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := ReadSigned(r, refreshCookieName, cfg.CookieKey)
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}
	err = cfg.Db.RevokeToken(r.Context(), token)
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}
	RespondWithJson(w, 204, struct{}{})
}

// get cookie handler
func (cfg *ApiConfig) ReaderCookieHandler(w http.ResponseWriter, r *http.Request) {
	value, err := ReadSigned(r, accessCookieName, cfg.CookieKey)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	type responseSturct struct {
		Value string `json:"value"`
	}
	RespondWithJson(w, 200, responseSturct{
		Value: value,
	})
}


func (cfg *ApiConfig) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request){
	userId, err := GetUserIdFromContext(r.Context())
	if err != nil {
		RespondWithError(w, 403, err.Error())
		return
	}
	user, err := cfg.Db.GetUserById(r.Context(), userId)
	if err != nil {
		RespondWithError(w, 400, err.Error())
	}
	type resBody struct{
		Email string `json:"email"`
	}
	res := resBody{
		Email: user.Email,
	}
	RespondWithJson(w, 200, res)
}


func (cfg *ApiConfig) LogoutHandler(w http.ResponseWriter, r *http.Request){
	token, err := ReadSigned(r, refreshCookieName, cfg.CookieKey)
	if err != nil {
		RespondWithError(w, 403, err.Error())
		return 
	}
	accesCookie := http.Cookie{
		Name:     accessCookieName,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
	}
	err = cfg.Db.RevokeToken(r.Context(), token)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	refreshCookie := http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
	}
	err = WriteSigned(w, accesCookie, cfg.CookieKey)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	err = WriteSigned(w, refreshCookie, cfg.CookieKey)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	RespondWithJson(w, 204, struct{}{})
}