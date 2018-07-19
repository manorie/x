package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var globalRouter *mux.Router

func init() {
	globalRouter = mux.NewRouter()
	initRoutes()
}

func initRoutes() {
	log.Println("initializing routes")

	globalRouter.HandleFunc("/api/addChecker", AddCheckerHandler)
	globalRouter.HandleFunc("/api/stopChecker", AddCheckerHandler)
	globalRouter.HandleFunc("/api/runChecker", AddCheckerHandler)
	globalRouter.HandleFunc("/api/deleteChecker", AddCheckerHandler)
	globalRouter.HandleFunc("/api/updateChecker", AddCheckerHandler)
	globalRouter.HandleFunc("/api/checksSince", AddCheckerHandler)
	http.Handle("/", globalRouter)
}

// AddCheckerHandler ...
func AddCheckerHandler(w http.ResponseWriter, r *http.Request) {
	ep, err := newEndPoint("https://ounass.ae/asdasda/asdads", "GET", true)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := newChecker(ep, 60*time.Second, 60*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	resp, _ := json.Marshal(ch)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

type CheckerResponse struct {
	Ch   *Checker `json:"checker"`
	Chks []*Check `json:"checks"`
}

type CheckersResponse struct {
	Chrs []*CheckerResponse `json:"checks"`
}

type BasicResponse struct {
	Msg string `json:"msg"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}
