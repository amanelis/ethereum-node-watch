package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

func main() {
    var router = mux.NewRouter()
    router.HandleFunc("/", index).Methods("GET")
    router.HandleFunc("/healthcheck", healthCheck).Methods("GET")
	router.HandleFunc("/message", handleQryMessage).Methods("GET")
    router.HandleFunc("/m/{msg}", handleUrlMessage).Methods("GET")

    fmt.Println("Server started at http://localhost:3000")
    log.Fatal(http.ListenAndServe(":3000", router))
}

func index(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]string{"message": "OK"})
}

func handleQryMessage(w http.ResponseWriter, r *http.Request) {
    vars := r.URL.Query()
    message := vars.Get("msg")

    json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func handleUrlMessage(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    message := vars["msg"]

    json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode("OK")
}
