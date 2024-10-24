package main

import (
	"fmt"
	"html/template"
	"lab0302/knn"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var recsys *knn.Knn

func piwaFunc(w http.ResponseWriter, r *http.Request) {
	// Retrieve 10 random beers from the system
	tenbeers := recsys.Get10RandomBeers()

	// Parse and execute the template with the 10 random beers
	tmpl, err := template.ParseFiles("pages/beer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, tenbeers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rekoFunc(w http.ResponseWriter, r *http.Request) {
	// gdy metoda inna niż POST - error 403
	if r.Method != "POST" {
		http.Error(w, "Only POST supported", http.StatusForbidden)
		return
	}
	// parsowanie danych z formularza
	r.ParseForm()
	// iteracja danych z formularza i jednoczesna ocena piw
	for name, element := range r.PostForm {
		rate, _ := strconv.ParseFloat(element[0], 64) // parsowanie
		id, _ := strconv.Atoi(name[4:])
		// wycięcie id
		beer := recsys.GetBeerByID(id)
		// pobranie piwa
		beer.Rate = rate
		// ocena piwa
		fmt.Println("Name:", name, ", Id:", id, ", Rate:", rate)
	}
	// pobranie rekomendacji
	recbeers := recsys.GetRecommendation()
	tmpl, _ := template.ParseFiles("pages/reko.html")
	tmpl.Execute(w, recbeers)
}

func main() {
	// These lines should remain
	rand.Seed(time.Now().UnixMilli())
	recsys = knn.Initialize()

	// Comment out or remove the code below
	// ...

	// Set up the HTTP server
	http.HandleFunc("/piwa", piwaFunc)
	http.HandleFunc("/reko", rekoFunc)
	fmt.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
