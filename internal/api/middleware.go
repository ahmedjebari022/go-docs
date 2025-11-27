package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ahmedjebari022/go-docs/internal/auth"
	"github.com/google/uuid"
)

type contextKey string
const (
	k = contextKey("userID")
)
func (cfg *ApiConfig) AuthMiddleware (next http.Handler) http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		value, err := ReadSigned(r,accessCookieName,cfg.CookieKey)
		if err != nil {
			respondWithError(w,401,"authentication error")
			return
		}
		userId, err := auth.ValidateJwt(cfg.SecretKey, value)
		if err != nil {
			respondWithError(w,401,"authentication error")
			return
		}
	
		ctx := context.WithValue(r.Context(),k,userId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w,r)
	})  
} 


func GetUserIdFromContext(ctx context.Context) (uuid.UUID, error) {
	if userId := ctx.Value(k).(uuid.UUID); userId != uuid.Nil{
		return userId, nil
	}
	return uuid.Nil, fmt.Errorf("couldn't get user id from context")
}