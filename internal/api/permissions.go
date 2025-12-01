package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ahmedjebari022/go-docs/internal/database"
	"github.com/google/uuid"
)

const (
	ViewerRole = "viewer"
	EditorRole = "editor"
	OwnerRole = "owner"
)

func  getDocumentAndUserFromUrl(r *http.Request) (userId, documentId uuid.UUID, err error) {
	documentIdString := r.PathValue("documentId")
	documentId, err = uuid.Parse(documentIdString)
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("400: document error")
	}
	userId, err = GetUserIdFromContext(r.Context())
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("401: not authenticated")
	}
	return userId, documentId, err
}

func (cfg *ApiConfig) checkIsUserOwner(ctx context.Context, user_id, document_id uuid.UUID) (bool, error){
	 ownerId, err := cfg.Db.GetDocumentOwnerId(ctx, document_id)
	 if err != nil {
		return false, err
	 }
	 return ownerId == user_id, nil
}

func (cfg *ApiConfig) requireOwnerShip(r *http.Request) (userId, documentId uuid.UUID, err error){
	userId, documentId, err = getDocumentAndUserFromUrl(r)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	isOwner, err := cfg.checkIsUserOwner(r.Context(), userId, documentId)
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("404: document not found")
	}
	if !isOwner {
		return uuid.Nil, uuid.Nil, fmt.Errorf("403: not authorized")
	}
	return userId, documentId, nil
}
func parseStatusFromError(err error) int {
	msg := err.Error()
	if len(msg) < 4{
		return 500
	} 
	statusCodeString := msg[:3]
	statusCode, statusErr := strconv.Atoi(statusCodeString)
	if statusErr != nil  {
		return 500
	}
	return statusCode
}

func (cfg *ApiConfig) AddCollaboratorToDocumentHandler(w http.ResponseWriter, r *http.Request){
	type requestBody struct{
		UserId 		uuid.UUID `json:"user_id"`
		Role 		string  `json:"role"`
	}

	_, documentId, err := cfg.requireOwnerShip(r)
	if err != nil {
		statusCode := parseStatusFromError(err)
		respondWithError(w,statusCode, err.Error() )
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

	
	_, documentId, err := cfg.requireOwnerShip(r) 
	if err != nil {
		statusCode := parseStatusFromError(err)
		respondWithError(w, statusCode, err.Error())
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
		respondWithError(w, 400, "invalid role")
		return
	}
	err = cfg.Db.UpdatePermission(r.Context(), database.UpdatePermissionParams{
		DocumentID: documentId,
		UserID: params.Id,
		Role: params.Role,
	})
	if err != nil {
		respondWithError(w, 400, err.Error())
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
	userId, documentId, err := getDocumentAndUserFromUrl(r)
	if err != nil {
		statusCode := parseStatusFromError(err)
		respondWithError(w, statusCode, err.Error() )
		return
	}

	u, err := cfg.Db.GetUsersFromDocument(r.Context(),documentId)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	
	owner, err := cfg.Db.GetDocumentOwner(r.Context(), documentId)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return 
	}

	if isCollaborator := userIsCollaborator(u, userId) || userId == owner.ID  ; !isCollaborator{
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
	
	ownerUr:= userRole{
		Email: owner.Email,
		Role: OwnerRole,
	}
	res.UserRoles = append(res.UserRoles, ownerUr)
	RespondWithJson(w, 200, res)
}

func (cfg *ApiConfig) DeleteUserFromCollaboration(w http.ResponseWriter, r *http.Request){
	type requestBody struct {
		Id 	uuid.UUID `json:"id"`
	}
	
	_, documentId, err := cfg.requireOwnerShip(r)
	if err != nil {
		statusCode := parseStatusFromError(err)
		respondWithError(w, statusCode, err.Error())
		return
	}

	var params requestBody
	decoder:= json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 400, err.Error())
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
	RespondWithJson(w, 204 ,struct{}{})
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