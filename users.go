package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	type reqparameters struct {
		Email string `json:"email"`
	}

	type respparameters struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
		Error      string    `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	reqParams := reqparameters{}
	respBody := respparameters{}
	statusCode := 201

	err := decoder.Decode(&reqParams)
	// bad request
	if err != nil {
		respBody.Error = fmt.Sprintf("%s", err)
		statusCode = 400
	} else if len(reqParams.Email) < 3 {
		respBody.Error = "Wrong/No Email"
		statusCode = 400
	} else {
		// good request
		User, err := Config.dbQueries.CreateUser(r.Context(), reqParams.Email)
		if err != nil {
			respBody.Error = fmt.Sprintf("Error creating user %s:%s", reqParams.Email, err)
			statusCode = 400
		} else {
			respBody.Id = User.ID
			respBody.Created_at = User.CreatedAt
			respBody.Updated_at = User.UpdatedAt
			respBody.Email = User.Email
		}
	}
	respJson, err := json.Marshal(respBody)
	if err != nil {
		fmt.Printf("Error marshalling json %s", err)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(respJson)
}
