package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"lab0302/knn"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var recsys *knn.Knn

func piwaJsonFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tenbeers := recsys.Get10RandomBeers()

	data, err := json.Marshal(tenbeers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func piwaFunc(w http.ResponseWriter, r *http.Request) {
	tenbeers := recsys.Get10RandomBeers()

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
	if r.Method != "POST" {
		http.Error(w, "Only POST supported", http.StatusForbidden)
		return
	}
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

func piwoJsonFunc(w http.ResponseWriter, r *http.Request) {
	idStr := strings.Trim(r.URL.Path, "/pid")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Nieprawidłowe id piwa", http.StatusBadRequest)
		return
	}

	beer := recsys.GetBeerByID(id)
	if beer == nil {
		http.Error(w, "Nie znaleziono piwa", http.StatusNotFound)
		return
	}

	similarBeers := recsys.GetThreeMostSimilarBeers(beer)

	// Prepare the response data
	response := struct {
		Beer         *knn.Beer   `json:"beer"`
		SimilarBeers []*knn.Beer `json:"similar_beers"`
	}{
		Beer:         beer,
		SimilarBeers: similarBeers,
	}

	w.Header().Set("Content-Type", "application/json")

	// Marshal the response data to JSON
	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func main() {
	rand.Seed(time.Now().UnixMilli())
	recsys = knn.Initialize()

	http.HandleFunc("/piwa", piwaFunc)
	http.HandleFunc("/piwajson", piwaJsonFunc)
	http.HandleFunc("/reko", rekoFunc)
	http.HandleFunc("/pid/", piwoJsonFunc)
	fmt.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
