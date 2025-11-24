package main

import (
	"encoding/json"
	"fmt"
	"internal/auth"
	"internal/database"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	type reqparameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type respparameters struct {
		Id            uuid.UUID `json:"id"`
		Created_at    time.Time `json:"created_at"`
		Updated_at    time.Time `json:"updated_at"`
		Email         string    `json:"email"`
		Error         string    `json:"error"`
		Is_Chirpy_Red bool      `json:"is_chirpy_red"`
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
				respBody.Is_Chirpy_Red = User.IsChirpyRed
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type respparameters struct {
		Id            uuid.UUID `json:"id"`
		Created_at    time.Time `json:"created_at"`
		Updated_at    time.Time `json:"updated_at"`
		Email         string    `json:"email"`
		Error         string    `json:"error"`
		Token         string    `json:"token"`
		Refresh_Token string    `json:"refresh_token"`
		Is_Chirpy_Red bool      `json:"is_chirpy_red"`
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
				tokenstring, _ := auth.MakeJWT(User.ID, Config.secret, time.Duration((3600 * int(time.Second))))
				refresh_token, _ := auth.MakeRefreshToken()
				Config.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{Token: refresh_token, UserID: User.ID, ExpiresAt: time.Now().AddDate(0, 0, 60)})

				respBody.Id = User.ID
				respBody.Created_at = User.CreatedAt
				respBody.Updated_at = User.UpdatedAt
				respBody.Email = User.Email
				respBody.Token = tokenstring
				respBody.Refresh_Token = refresh_token
				respBody.Is_Chirpy_Red = User.IsChirpyRed
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

func RefreshToken(w http.ResponseWriter, r *http.Request) {

	type respparameters struct {
		Token string `json:"token"`
	}
	respBody := respparameters{}
	statusCode := 201

	refreshToken := r.Header.Get("Authorization")
	refreshToken, _ = strings.CutPrefix(refreshToken, "Bearer ")
	dbToken, err := Config.dbQueries.GetUserFromRefreshToken(r.Context(), refreshToken)

	if err != nil || dbToken.RevokedAt.Valid {
		statusCode = 401
	} else {
		statusCode = 200
		token, _ := auth.MakeJWT(dbToken.UserID, Config.secret, time.Hour)
		respBody.Token = token
	}
	respJson, _ := json.Marshal(respBody)
	w.WriteHeader(statusCode)
	w.Write(respJson)
}

func RevokeToken(w http.ResponseWriter, r *http.Request) {

	statusCode := 201

	refreshToken := r.Header.Get("Authorization")
	refreshToken, _ = strings.CutPrefix(refreshToken, "Bearer ")

	err := Config.dbQueries.RevokeRefreshToken(r.Context(), refreshToken)

	if err != nil {
		statusCode = 401
	} else {
		statusCode = 204
	}

	w.WriteHeader(statusCode)

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	type reqparameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type respparameters struct {
		Id            uuid.UUID `json:"id"`
		Created_at    time.Time `json:"created_at"`
		Updated_at    time.Time `json:"updated_at"`
		Email         string    `json:"email"`
		Error         string    `json:"error"`
		Is_Chirpy_Red bool      `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(r.Body)
	reqParams := reqparameters{}
	respBody := respparameters{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	accessToken := r.Header.Get("Authorization")
	accessToken, _ = strings.CutPrefix(accessToken, "Bearer ")

	User, err := auth.ValidateJWT(accessToken, Config.secret)

	if err != nil || reqParams.Email == "" || reqParams.Password == "" {
		w.WriteHeader(401)
		return
	} else {
		hash, err := auth.HashPassword(reqParams.Password)
		if err != nil {
			w.WriteHeader(401)
			return
		}
		dbUser, err := Config.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{ID: User, Email: reqParams.Email, Password: hash})
		if err != nil {
			w.WriteHeader(401)
			return
		}
		respBody.Id = dbUser.ID
		respBody.Created_at = dbUser.CreatedAt
		respBody.Updated_at = dbUser.UpdatedAt
		respBody.Email = dbUser.Email
		respBody.Is_Chirpy_Red = dbUser.IsChirpyRed
		respJson, _ := json.Marshal(respBody)
		w.WriteHeader(200)
		w.Write(respJson)
	}

}

func UpgradeUser(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != Config.polkakey {
		w.WriteHeader(401)
		return
	}

	type reqData struct {
		User_id uuid.UUID `json:"user_id"`
	}
	type reqparameters struct {
		Event string  `json:"event"`
		Data  reqData `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	reqParams := reqparameters{}
	err = decoder.Decode(&reqParams)
	if err != nil || reqParams.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	dbUser, err := Config.dbQueries.UpgradeUser_is_red(r.Context(), reqParams.Data.User_id)
	if err != nil || dbUser.Email == "" {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(204)
}
