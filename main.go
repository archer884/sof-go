package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/archer884/sof-go/cookies"
	"github.com/gorilla/mux"
)

type cookieResponse struct {
	Timestamp time.Time
	Cookie    *cookies.Cookie
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Provide quote directory")
		os.Exit(1)
	}

	service, err := cookies.New(args[0], rand.New(rand.NewSource(time.Now().UnixNano())))
	if err != nil {
		fmt.Println("Error creating cookie service")
		os.Exit(2)
	}

	var getCookie = func(w http.ResponseWriter, r *http.Request) {
		cookie := service.GetCookie()
		response := cookieResponse{time.Now().UTC(), &cookie}
		encoded, err := json.Marshal(response)
		if err != nil {
			fmt.Fprint(w, err)
		} else {
			fmt.Fprint(w, string(encoded))
		}
	}

	var getCookieWithCategory = func(w http.ResponseWriter, r *http.Request) {
		cookie, err := service.ByCategory(mux.Vars(r)["category"])
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		response := cookieResponse{time.Now().UTC(), &cookie}
		encoded, encErr := json.Marshal(response)
		if err != nil {
			fmt.Fprint(w, encErr)
		} else {
			fmt.Fprint(w, string(encoded))
		}
	}

	rtr := mux.NewRouter()
	rtr.HandleFunc("/api/cookie", getCookie).Methods("GET")
	rtr.HandleFunc("/api/cookie/{category}", getCookieWithCategory).Methods("GET")

	http.Handle("/", rtr)
	http.ListenAndServe(":5000", nil)
}
