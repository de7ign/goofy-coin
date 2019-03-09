package main

import (
	"log"
	"net/http"
	"time"
)

func reqLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request Time %s ", time.Now())
		log.Printf("Method %s ", r.Method)
		log.Printf("Request URI %s ", r.RequestURI)
		log.Printf("Remote address %s", r.RemoteAddr)
		next.ServeHTTP(w, r)
	}
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/index.html")
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/dashboard.html")
}

func main() {
	http.HandleFunc("/", reqLogger(indexHandler))
	http.HandleFunc("/dashboard", reqLogger(dashboardHandler))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./assets/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./assets/css"))))
	log.Printf("App running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
