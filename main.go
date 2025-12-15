package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movies struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var movies []Movies

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

func generateUniqueId() string {
	for {
		newId := strconv.Itoa(rand.Intn(1000000000))

		exists := false
		for _, movie := range movies {
			if movie.ID == newId {
				exists = true
				break
			}
		}
		if !exists {
			return newId
		}
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Получение параметров из URL запроса
	// params := map[string]string { "id": "123", }
	// Берет наш маршрут , сравнивает с фактичесим URL (прим. movies/123) и извлекаем значение
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			// Кодируем movies в JSON и отправляем в ответ
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Movie not found"})
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Movie not found"})
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movies
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}
	movie.ID = generateUniqueId()
	movies = append(movies, movie)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movies)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var updatedMovie Movies
	err := json.NewDecoder(r.Body).Decode(&updatedMovie)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
	}
	for index, item := range movies {
		if item.ID == params["id"] {
			updatedMovie.ID = item.ID

			if updatedMovie.Director == nil && item.Director != nil {
				updatedMovie.Director = item.Director
			}
			movies[index] = updatedMovie

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Movie not found"})
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movies{ID: "1", Title: "Iron man", Director: &Director{FirstName: "John", LastName: "Doe"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Server starting at port:8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
