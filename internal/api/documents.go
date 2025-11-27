package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ahmedjebari022/go-docs/internal/database"
	"github.com/google/uuid"
)


func (cfg *ApiConfig) CreateDocumentHandler(w http.ResponseWriter, r *http.Request){
		userId, err := GetUserIdFromContext(r.Context())
		if err != nil {
			respondWithError(w, 401, err.Error())
			return
		}

		type requestBody struct{
			Name 	string `json:"name"`
		}
		type responseBody struct{
			Name 	string `json:"name"`
			DocumentUrl 	string `json:"document_url"`
			OwnerId 			string `json:"owner_id"`
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithError(w, 400, err.Error())
			return
		}

		var params requestBody
		err = json.Unmarshal(data, &params)
		if err != nil  {
			respondWithError(w, 500,  err.Error())
			return
		}
		documentID := uuid.New()
		documentPath := filepath.Join(cfg.AssetsPath, documentID.String()) + ".json"
		err = EnsureDirExists(cfg.AssetsPath)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}
		_, err = os.Create(documentPath)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}
		documentUrl := fmt.Sprintf("http://localhost:%s/assets/%s.json", cfg.Port, documentID)
		document, err := cfg.Db.CreateDocument(r.Context(), database.CreateDocumentParams{
			ID: documentID,
			Name: params.Name,
			DocumentUrl: documentUrl,
			OwnerID: userId,
		})
		if err != nil {
			os.Remove(documentPath)
			respondWithError(w, 500, err.Error())
			return
		}
		RespondWithJson(w, 200, responseBody{
			Name: document.Name,
			DocumentUrl: document.DocumentUrl,
			OwnerId: document.OwnerID.String(),	
		})

}



func EnsureDirExists(pathdir string) error{
	_, err := os.Stat(pathdir)
	if os.IsNotExist(err){
		err = os.Mkdir(pathdir,0755)
		if err != nil {
			return err
		}
	}
	return nil
}