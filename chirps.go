package main

import (
	"encoding/json"
	"fmt"
	"internal/database"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func CreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body    string    `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	type respParams struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		User_id    uuid.UUID `json:"user_id"`
		Error      string    `json:"error"`
		Valid      bool      `json:"valid"`
	}

	respBody := respParams{}
	statusCode := 201
	// request errors
	if err != nil {
		respBody.Error = fmt.Sprintf("%s", err)
		respBody.Valid = false
		statusCode = 400
	} else if len(params.Body) > 140 || len(params.Body) == 0 {
		respBody.Error = "Chirp is too long/short"
		respBody.Valid = false
		statusCode = 400
	} else if len(params.User_id) == 0 {
		respBody.Error = "No user provided"
		respBody.Valid = false
		statusCode = 400
	} else {
		// request ok
		dbparams := database.CreateChirpParams{
			Body:   cleanChirp(params.Body),
			UserID: params.User_id,
		}
		Chirp, err := Config.dbQueries.CreateChirp(r.Context(), dbparams)
		if err != nil {
			respBody.Error = fmt.Sprintf("Error saving chirp:%s", err)
			statusCode = 400
		} else {
			respBody.Error = ""
			respBody.Valid = true
			statusCode = 201
			respBody.Id = Chirp.ID
			respBody.Created_at = Chirp.CreatedAt
			respBody.Updated_at = Chirp.UpdatedAt
			respBody.Body = Chirp.Body
			respBody.User_id = Chirp.UserID
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

func cleanChirp(dirty string) string {
	if len(dirty) == 0 {
		return ""
	}
	words := strings.Split(dirty, " ")
	for i, word := range words {
		lowWord := strings.ToLower(word)
		if lowWord == "kerfuffle" ||
			lowWord == "sharbert" ||
			lowWord == "fornax" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func GetChirps(w http.ResponseWriter, r *http.Request) {

	type respItem struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		User_id    uuid.UUID `json:"user_id"`
		Error      string    `json:"error"`
	}
	respBody := []respItem{}
	statusCode := 200

	Chirps, err := Config.dbQueries.GetChirps(r.Context())

	if err != nil {
		respBody = append(respBody, respItem{Error: fmt.Sprintf("Error getting chirps from db:%s", err)})
		statusCode = 400
	} else {
		for _, Chirp := range Chirps {
			respBody = append(respBody, respItem{Id: Chirp.ID,
				Created_at: Chirp.CreatedAt,
				Updated_at: Chirp.UpdatedAt,
				Body:       Chirp.Body,
				User_id:    Chirp.UserID,
			})
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

func GetChirp(w http.ResponseWriter, r *http.Request) {

	type respItem struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		User_id    uuid.UUID `json:"user_id"`
		Error      string    `json:"error"`
	}
	respBody := respItem{}
	statusCode := 200

	chirpId := r.PathValue("id")
	if chirpId == "" {
		respBody.Error = "No chirp ID provided"
		statusCode = 404
	} else {
		chirpUUID, err := uuid.Parse(chirpId)
		if err != nil {
			respBody.Error = "ID provided is not uuid"
			statusCode = 404
		} else {

			Chirp, err := Config.dbQueries.GetChirp(r.Context(), chirpUUID)
			if err != nil {
				respBody.Error = fmt.Sprintf("Error getting chirp from db:%s", err)
				statusCode = 404
			} else {
				respBody.Id = Chirp.ID
				respBody.Created_at = Chirp.CreatedAt
				respBody.Updated_at = Chirp.UpdatedAt
				respBody.Body = Chirp.Body
				respBody.User_id = Chirp.UserID
				statusCode = 200
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
