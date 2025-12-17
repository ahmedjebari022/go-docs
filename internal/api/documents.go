package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ahmedjebari022/go-docs/internal/database"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Document struct {
	Blocs []Bloc `json:"blocs"`
}
type Bloc struct {
	Type  string  `json:"type"`
	Children []Child `json:"children"`
}
type Child struct {
	Text  string  `json:"text"`
	Bold  bool     `json:"bold,omitempty"`
	Italic bool 	`json:"italic,omitempty"`
}

func (cfg *ApiConfig) CreateDocumentHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserIdFromContext(r.Context())
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}

	type requestBody struct {
		Name string `json:"name" validate:"required"`
	}
	type responseBody struct {
		Id      uuid.UUID `json:"document_id"`
		Name    string    `json:"document_name"`
		OwnerId string    `json:"owner_id"`
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}
	defer r.Body.Close()
	var params requestBody
	err = json.Unmarshal(data, &params)
	if err != nil {
		fmt.Printf("error while unmarhsalling %s \n",err.Error())
		RespondWithError(w, 500, err.Error())
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(params)
	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}
	documentID := uuid.New()
	err = EnsureDirExists(cfg.AssetsPath)
	if err != nil {
		fmt.Println("error while error while ensure dir exists")
		RespondWithError(w, 500, err.Error())
		return
	}
	documentPath := generatePathFromId(documentID.String(), cfg.AssetsPath)
	file, err := os.Create(documentPath)
	if err != nil {
		fmt.Println("error while creating folder")
		RespondWithError(w, 500, err.Error())
		return
	}
	defer file.Close()
	err = WriteToFile(file, []Bloc{})
	if err != nil {
		fmt.Println("error while writing to file")
		RespondWithError(w, 500, err.Error())
		return
	}
	document, err := cfg.Db.CreateDocument(r.Context(), database.CreateDocumentParams{
		ID:      documentID,
		Name:    params.Name,
		OwnerID: userId,
	})
	if err != nil {
		os.Remove(documentPath)
		RespondWithError(w, 500, err.Error())
		return
	}
	RespondWithJson(w, 201, responseBody{
		Id:      document.ID,
		Name:    document.Name,
		OwnerId: document.OwnerID.String(),
	})
}

func (cfg *ApiConfig) GetDocumentsByUserHandler(w http.ResponseWriter, r *http.Request) {

	type responseDocument struct {
		DocumentId   uuid.UUID `json:"document_id"`
		DocumentName string    `json:"document_name"`
	}
	type responseBody struct {
		Documents []responseDocument `json:"documents"`
	}

	userId, err := GetUserIdFromContext(r.Context())
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}

	documents, err := cfg.Db.GetDocumentsByUser(r.Context(), userId)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	var res responseBody
	for _, d := range documents {
		res.Documents = append(res.Documents, responseDocument{
			DocumentId:   d.ID,
			DocumentName: d.Name,
		})
	}
	RespondWithJson(w, 200, res)
}

func (cfg *ApiConfig) GetDocumentHandler(w http.ResponseWriter, r *http.Request) {
	
	type documentMetadata struct {
		DocumentId   uuid.UUID `json:"document_id"`
		DocumentName string    `json:"document_name"`
	}
	type resBody struct {
		MetaData  	documentMetadata `json:"meta_data"`
		DocumentContent []Bloc `json:"document_content"`
	}

	documentId := r.PathValue("documentId")
	userId, err := GetUserIdFromContext(r.Context())
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}
	if documentId == "" {
		RespondWithError(w, 400, "Missing document Id from the request")
		return
	}

	id, err := uuid.Parse(documentId)
	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}
	documentMetaData, err := cfg.Db.GetDocument(r.Context(), id)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	if documentMetaData.OwnerID != userId {
		_, err = cfg.Db.GetUserPermission(r.Context(), database.GetUserPermissionParams{
			UserID:     userId,
			DocumentID: id,
		})
		if err != nil {
			RespondWithError(w, 403, "Not Authorized to view this Document")
			return
		}
	}
	documentPath := generatePathFromId(documentId, cfg.AssetsPath)
	documentContent, err := ReadFromFile(documentPath)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	RespondWithJson(w, 200, resBody{
		MetaData: documentMetadata{
			DocumentId: documentMetaData.ID,
			DocumentName: documentMetaData.Name,
		},
		DocumentContent: documentContent,
	})
}

func (cfg *ApiConfig) UpdateDocumentHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserIdFromContext(r.Context())
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}
	documentId := r.PathValue("documentId")
	if documentId == "" {
		RespondWithError(w, 400, "Missing Document id from query")
		return
	}
	id, err := uuid.Parse(documentId)
	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}

	ownerId, err := cfg.Db.GetDocumentOwnerId(r.Context(), id)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	if ownerId != userId {
		role, err := cfg.Db.GetUserPermission(r.Context(), database.GetUserPermissionParams{
			UserID:     userId,
			DocumentID: id,
		})
		if err != nil || role != EditorRole {
			RespondWithError(w, 403, "user not authorized")
			return
		}
	}
	sizeLimit := 1 << 20
	r.Body = http.MaxBytesReader(w, r.Body, int64(sizeLimit))
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var params []Bloc
	if err := decoder.Decode(&params); err != nil {
		RespondWithError(w, 400, fmt.Sprintf("ici %s",err.Error()))
		return
	}
	if len(params) > 50000 {
		RespondWithError(w, 413, "document too large")
		return
	}
	ctx, err := cfg.DbC.Begin()
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	defer ctx.Rollback()
	qtx := cfg.Db.WithTx(ctx)
	err = qtx.UpdateDocument(r.Context(), id)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	path := generatePathFromId(documentId, cfg.AssetsPath)
	tmp, err := os.CreateTemp(cfg.AssetsPath, "tmp.json")
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	defer tmp.Close()
	err = WriteToFile(tmp, params)
	if err != nil {
		os.Remove(tmp.Name())
		RespondWithError(w, 500, err.Error())
		return
	}
	err = ctx.Commit()
	if err != nil {
		os.Remove(tmp.Name())
		RespondWithError(w, 500, err.Error())
		return
	}
	err = os.Rename(tmp.Name(), path)
	if err != nil {
		os.Remove(tmp.Name())
		RespondWithError(w, 500, err.Error())
		return
	}
	RespondWithJson(w, 204, struct{}{})
}

func (cfg *ApiConfig) DeleteDocumentHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserIdFromContext(r.Context())
	if err != nil {
		RespondWithError(w, 401, "action require authentication")
		return
	}
	documentIdString := r.PathValue("documentId")
	fmt.Printf("documentId: %s\n",documentIdString)
	documentId, err := uuid.Parse(documentIdString)
	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}

	ownerId, err := cfg.Db.GetDocumentOwnerId(r.Context(), documentId)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	if ownerId != userId {
		RespondWithError(w, 403, "unauthorized")
		return
	}
	ctx, err := cfg.DbC.Begin()
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	tx := cfg.Db.WithTx(ctx)
	defer ctx.Rollback()

	err = tx.DeleteDocument(r.Context(), documentId)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	err = ctx.Commit()
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	path := generatePathFromId(documentIdString, cfg.AssetsPath)
	err = os.Remove(path)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	RespondWithJson(w, 204, struct{}{})
}

func generatePathFromId(id, assets string) string {
	path := filepath.Join(assets, id)
	return path + ".json"
}

func EnsureDirExists(pathdir string) error {
	_, err := os.Stat(pathdir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(pathdir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteToFile(file *os.File, documentContent []Bloc) error {
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(documentContent); err != nil {
		return err
	}
	return nil
}

func ReadFromFile(filepath string) ([]Bloc, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var document []Bloc
	if err := decoder.Decode(&document); err != nil {
		return nil, err
	}
	return document, nil
}
