package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type HTMLDir struct {
	Dir http.Dir
}

func (d HTMLDir) Open(name string) (http.File, error) {
	f, err := d.Dir.Open(name)

	if os.IsNotExist(err) {
		if f, err := d.Dir.Open("index.html"); err == nil {
			return f, nil
		}
	}

	return f, err
}

type AppHandler struct {
	Handler func(context *Context, w http.ResponseWriter, r *http.Request)
	Context *Context
}

func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ah.Handler(ah.Context, w, r)
}

func IndexHandler(context *Context, w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func ProtectedRoute(context *Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Protected")
}

func main() {
	router := mux.NewRouter()
	router.Use(LoggingMiddleware)

	config := NewConfigFromEnv()
	db, err := NewDB(config.DB)

	if err != nil {
		log.Fatal("Failed to connect to the database")
	}

	db.CreateSchema()

	context := &Context{config, db}

	api := router.PathPrefix("/api/").Subrouter()
	api.Handle("/create-account", AppHandler{CreateUserHandler, context}).Methods("POST")
	api.Handle("/login", AppHandler{LoginUserHandler, context}).Methods("POST")

	router.Handle("/", AppHandler{IndexHandler, context}).Methods("GET")

	protected := router.PathPrefix("/protected").Subrouter()
	protected.Handle("/path", AppHandler{ProtectedRoute, context}).Methods("GET")
	protected.Use(verifyJWT(config.JWTSecret))

	fs := http.FileServer(HTMLDir{http.Dir(config.StaticDir)})
	router.Handle("/{path:.*}", http.StripPrefix("/", fs))

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
