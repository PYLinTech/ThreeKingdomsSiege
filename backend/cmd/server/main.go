package main

import (
	"log"
	"net/http"
	"os"
	"time"

	auth "fangguo/backend/internal/auth"
)

func main() {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./data/app.db"
	}
	if err := os.MkdirAll("./data", 0o750); err != nil {
		log.Fatal(err)
	}
	store, err := auth.OpenSQLite(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	cache := auth.NewCache()
	issuer, err := auth.NewIssuer(cache, store, 24*time.Hour)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	authHandler := auth.AuthHandler{Store: store, Issuer: issuer}
	mux.Handle("/api/auth/register", authHandler)
	mux.Handle("/api/auth/login", authHandler)
	mux.Handle("/api/auth/token", auth.Handler{Issuer: issuer})
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("token service listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
