package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type DateResp struct {
	Date  string `json:"date"`
	Login string `json:"login"`
}

const login = "levchik"

// хэндлер для даты
func dateHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	ddmmyy := now.Format("020106")
	path := r.URL.Path[1:]

	// Проверяем, совпадает ли путь с сегодняшней датой
	if path == ddmmyy+"/" {
		w.Header().Set("Content-Type", "application/json")
		resp := DateResp{
			Date:  now.Format("02-01-2006"), // DD-MM-YYYY
			Login: login,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	http.NotFound(w, r)
}

func reverseHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	word := parts[3]

	match, _ := regexp.MatchString("^[a-z]+$", word)
	if !match {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	runes := []rune(word)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(string(runes)))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", dateHandler)

	http.HandleFunc("/api/rv/", reverseHandler)

	log.Printf("Listening on :%s ...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
