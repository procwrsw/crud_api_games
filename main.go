package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Game struct {
	ID        string `json:"id"`
	Isbn      string `json:"isbn"`
	Name      string `json:"name"`
	Developer *Developer
}

type Developer struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var games []Game

func getGames(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

func getGame(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for _, item := range games {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func updateGame(w http.ResponseWriter, req *http.Request) {
	// at first we want to delete game
	// next step is to create a game
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for index, item := range games {
		if item.ID == params["id"] {
			games = append(games[:index], games[index+1:]...)
			var game Game
			_ = json.NewDecoder(req.Body).Decode(&game)
			game.ID = strconv.Itoa(rand.Intn(100))
			games = append(games, game)
			json.NewEncoder(w).Encode(games)
		}
	}

}

func createGame(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var game Game
	_ = json.NewDecoder(req.Body).Decode(&game)
	game.ID = strconv.Itoa(rand.Intn(100))
	games = append(games, game)
	json.NewEncoder(w).Encode(games)
}

func deleteGame(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for index, item := range games {
		if item.ID == params["id"] {
			games = append(games[:index], games[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(games)
}

func main() {
	r := mux.NewRouter()

	games = append(games, Game{ID: "1", Isbn: "334657", Name: "Game one", Developer: &Developer{Firstname: "Adam", Lastname: "Sandler"}})
	games = append(games, Game{ID: "2", Isbn: "456987", Name: "Game two", Developer: &Developer{Firstname: "Mark", Lastname: "Thomson"}})

	r.HandleFunc("/games", getGames).Methods("GET")
	r.HandleFunc("/games/{id}", getGame).Methods("GET")
	r.HandleFunc("/games", createGame).Methods("POST")
	r.HandleFunc("/games/{id}", updateGame).Methods("PUT")
	r.HandleFunc("/games/{id}", deleteGame).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
