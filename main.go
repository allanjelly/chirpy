package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiconfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	secret         string
	polkakey       string
}

var Config apiconfig

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Print("Error connecting to db")
		os.Exit(1)
	}
	Config.polkakey = os.Getenv("POLKA_KEY")
	Config.secret = os.Getenv("SECRET")
	Config.dbQueries = database.New(db)

	ServeMux := http.NewServeMux()
	//static files handler
	ServeMux.Handle("/app/", http.StripPrefix("/app", Config.middlewareMetricsInc(http.FileServer(http.Dir("./static")))))

	//admin
	ServeMux.HandleFunc("GET /api/healthz", WatchdogHandler)
	ServeMux.HandleFunc("GET /admin/metrics", MetricsHandler)
	ServeMux.HandleFunc("POST /admin/reset", MetricsResetHandler)

	//api
	ServeMux.HandleFunc("POST /api/users", CreateUser)
	ServeMux.HandleFunc("PUT /api/users", UpdateUser)
	ServeMux.HandleFunc("POST /api/login", UserLogin)
	ServeMux.HandleFunc("POST /api/refresh", RefreshToken)
	ServeMux.HandleFunc("POST /api/revoke", RevokeToken)
	ServeMux.HandleFunc("POST /api/polka/webhooks", UpgradeUser)

	ServeMux.HandleFunc("POST /api/chirps", CreateChirp)
	ServeMux.HandleFunc("GET /api/chirps", GetChirps)
	ServeMux.HandleFunc("GET /api/chirps/{id}", GetChirp)
	ServeMux.HandleFunc("DELETE /api/chirps/{chirpID}", DeleteChirp)

	Config.fileserverHits.Store(0)
	var Server http.Server
	Server.Handler = ServeMux
	Server.Addr = ":8080"

	Server.ListenAndServe()

}
