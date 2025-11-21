package main

import (
	"encoding/json"
	"fmt"
	"internal/auth"
	"internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	type reqparameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
	} else if len(reqParams.Password) < 5 {
		respBody.Error = "No/too short (min.5) password"
		statusCode = 400
	} else {
		// good request
		hash, err := auth.HashPassword(reqParams.Password)
		if err != nil {
			respBody.Error = "Could not hash password"
			statusCode = 400
		} else {
			User, err := Config.dbQueries.CreateUser(r.Context(), database.CreateUserParams{Email: reqParams.Email, Password: hash})
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
	}
	respJson, err := json.Marshal(respBody)
	if err != nil {
		fmt.Printf("Error marshalling json %s", err)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(respJson)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {

	type reqparameters struct {
		Email              string `json:"email"`
		Password           string `json:"password"`
		Expires_in_seconds int    `json:"expires_in_seconds"`
	}

	type respparameters struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
		Error      string    `json:"error"`
		Token      string    `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	reqParams := reqparameters{}
	respBody := respparameters{}
	statusCode := 201

	err := decoder.Decode(&reqParams)
	// bad request
	if err != nil {
		respBody.Error = fmt.Sprintf("%s", err)
		statusCode = 401
	} else if len(reqParams.Email) < 3 {
		respBody.Error = "Wrong/No Email"
		statusCode = 401
	} else if len(reqParams.Password) < 5 {
		respBody.Error = "No/too short (min.5) password"
		statusCode = 401
	} else {
		// good request
		User, err := Config.dbQueries.GetUser(r.Context(), reqParams.Email)
		if err != nil {
			respBody.Error = "Wrong user/password"
			statusCode = 401
		} else {
			ok, err := auth.CheckPasswordHash(reqParams.Password, User.Password)
			//login ok
			if ok && err == nil {
				if !(reqParams.Expires_in_seconds > 0 && reqParams.Expires_in_seconds < 3600) {
					reqParams.Expires_in_seconds = 3600
				}
				tokenstring, _ := auth.MakeJWT(User.ID, Config.secret, time.Duration((reqParams.Expires_in_seconds * int(time.Second))))
				respBody.Id = User.ID
				respBody.Created_at = User.CreatedAt
				respBody.Updated_at = User.UpdatedAt
				respBody.Email = User.Email
				respBody.Token = tokenstring
				statusCode = 200
			} else {
				//wrong pass
				statusCode = 401
				respBody.Error = "Wrong user/password"
			}
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
