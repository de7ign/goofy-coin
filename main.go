package main

import (
	"log"
	"net/http"

)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/index.html")
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/dashboard.html")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
	log.Printf("App running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}