package api

import (
	"encoding/json"
	"net/http"
)


func RespondWithJson(w http.ResponseWriter, code int, payload any)error{
	w.Header().Set("Content-Type","application/json; charset=utf-8")
	res, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.WriteHeader(code)
	w.Write(res)
	return nil	
}


func respondWithError(w http.ResponseWriter, code int, msg string)error{
	return RespondWithJson(w,code,map[string]string{"error":msg})	
}

