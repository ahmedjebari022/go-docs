package api

import (
	"encoding/json"
	"net/http"

	"github.com/ahmedjebari022/go-docs/internal/database"
	"github.com/google/uuid"
)

const (
	ViewerRole = "viewer"
	EditorRole = "editor"
	OwnerRole = "owner"
)


func (cfg *ApiConfig) checkIfUserIsOwner(r *http.Request, userId, documentId uuid.UUID)(bool, error){
	ownerId, err := cfg.Db.GetDocumentOwner(r.Context(), documentId)
	if err != nil {
		return false, err
	}
	return ownerId == userId, nil
}

func (cfg *ApiConfig) AddCollaboratorToDocumentHandler(w http.ResponseWriter, r *http.Request){
	type requestBody struct{
		UserId 		uuid.UUID `json:"user_id"`
		Role 		string  `json:"role"`
	}

	documentIdString := r.URL.Query().Get("documentId")
	documentId, err := uuid.Parse(documentIdString)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return 
	}

	userId, err := GetUserIdFromContext(r.Context())
	if err != nil {
		respondWithError(w, 401, err.Error())
		return 
	}

	isOwner, err := cfg.checkIfUserIsOwner(r, userId, documentId)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	if !isOwner{
		respondWithError(w, 403, "not authorized")
		return
	}
	var params requestBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 400, err.Error())
		return 
	}
	defer r.Body.Close()
	if params.Role != EditorRole && params.Role != ViewerRole{
		respondWithError(w, 400, "wrong role value")
		return
	}		

	err = cfg.Db.CreatePermission(r.Context(), database.CreatePermissionParams{
		UserID: params.UserId,
		DocumentID: documentId,
		Role: params.Role,

	})

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	RespondWithJson(w, http.StatusCreated, struct{}{})
}


func (cfg *ApiConfig) UpdateUserPermissionHandler(w http.ResponseWriter, r *http.Request){
	
	type requestBody struct{
		Id 		uuid.UUID `json:"id"`
		Role 	string `json:"role"`
	}

	documentIdString := r.URL.Query().Get("documentId")
	documentId, err := uuid.Parse(documentIdString)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return 
	}
	userId, err := GetUserIdFromContext(r.Context())
	if err != nil {
		respondWithError(w, 403, err.Error())
		return
	}
	isOwner, err := cfg.checkIfUserIsOwner(r, userId, documentId)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return 
	}
	if !isOwner {
		respondWithError(w, 403, "not authorized")
		return
	}
	var params requestBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	defer r.Body.Close()
	if params.Role != EditorRole && params.Role != ViewerRole{
		respondWithError(w, 500, "invalid role")
		return
	}
	err = cfg.Db.UpdatePermission(r.Context(), database.UpdatePermissionParams{
		DocumentID: documentId,
		UserID: params.Id,
		Role: params.Role,
	})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	RespondWithJson(w, 200, struct{}{})
}


func (cfg *ApiConfig) GetCollaboratorsHandler (w http.ResponseWriter, r *http.Request){
	type userRole struct {
		Email 	string `json:"email"`
		Role 	string `json:"role"`
	}
	type responseBody struct{
		UserRoles	 []userRole `json:"userRoles"`
	}
	documentIdString := r.URL.Query().Get("documentId")
	documentId, err := uuid.Parse(documentIdString)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	userId, err := GetUserIdFromContext(r.Context())
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}
	u, err := cfg.Db.GetUsersFromDocument(r.Context(),documentId)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	
	ownerId, err := cfg.Db.GetDocumentOwner(r.Context(), documentId)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return 
	}

	if isCollaborator := userIsCollaborator(u, userId) || userId == ownerId  ; !isCollaborator{
		respondWithError(w, 403, "not authorized")
		return
	}
	res := responseBody{}
	for _, v := range u {
		ur := userRole{
			Email: v.Email,
			Role: v.Role,
		}
		res.UserRoles = append(res.UserRoles, ur)
	}

	owner, err := cfg.Db.GetUserById(r.Context(), ownerId)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	
	ownerEr:= userRole{
		Email: owner.Email,
		Role: OwnerRole,
	}
	res.UserRoles = append(res.UserRoles, ownerEr)
	RespondWithJson(w, 200, res)
}

func (cfg *ApiConfig) DeleteUserFromCollaboration(w http.ResponseWriter, r *http.Request){
	type requestBody struct {
		Id 	uuid.UUID `json:"id"`
	}
	documentIdString := r.URL.Query().Get("documentId")
	documentId, err := uuid.Parse(documentIdString)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	
	userId, err := GetUserIdFromContext(r.Context())
	if err != nil {
		respondWithError(w, 401, err.Error())
		return 
	}

	isOwner, err := cfg.checkIfUserIsOwner(r, userId, documentId)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	if !isOwner {
		respondWithError(w, 403, "not authorized")
		return
	} 
	var params requestBody
	docoder := json.NewDecoder(r.Body)
	if err := docoder.Decode(&params); err != nil {
		respondWithError(w, 500, err.Error())
		return 
	}
	defer r.Body.Close()

	err = cfg.Db.DeletePermission(r.Context(), database.DeletePermissionParams{
		UserID: params.Id,
		DocumentID: documentId,
	})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	RespondWithJson(w, 200 , map[string]string{
		"message": "permission Revoked succesfully",
	})

}


func userIsCollaborator(u []database.GetUsersFromDocumentRow, userId uuid.UUID) bool {
	isCollaborator := false
	for _, v := range u {
		if v.ID == userId{
			isCollaborator = true
			break
		}
	}
	return isCollaborator
}