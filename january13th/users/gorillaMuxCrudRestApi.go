package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     int    `json:"age"`
}

var users []User
var idCounter int

func main() {
	setupAndRunUsersServer()
}

func setupAndRunUsersServer() {
	log.Println("Initializing server...")
	router := mux.NewRouter()
	log.Println("...done")

	setupDatabase()

	log.Println("Configuring server...")
	router.HandleFunc("/api/users", getUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/api/users", createUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/api/users/{id}", deleteUser).Methods("DELETE")
	log.Println("...done")

	log.Println("Running server...")
	log.Fatal(http.ListenAndServe(":7080", router))
	log.Println("...done")
}

func setupDatabase() {
	users = append(users, User{ID: 1, Name: "Eren", Surname: "Jeager", Age: 17})
	users = append(users, User{ID: 2, Name: "Mikasa", Surname: "Ackerman", Age: 17})
	users = append(users, User{ID: 3, Name: "Armin", Surname: "Arlert", Age: 16})
	users = append(users, User{ID: 4, Name: "Levi", Surname: "Ackerman", Age: 27})
	users = append(users, User{ID: 5, Name: "Annie", Surname: "Leonheart", Age: 18})
	users = append(users, User{ID: 6, Name: "Frida", Surname: "Reiss", Age: 15})
	users = append(users, User{ID: 7, Name: "Sasha", Surname: "Brouse", Age: 16})
	users = append(users, User{ID: 8, Name: "Zoe", Surname: "Hange", Age: 29})
	idCounter = 8
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	idCounter++
	user.ID = idCounter

	users = append(users, user)

	json.NewEncoder(w).Encode(user)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for _, user := range users {
		id, _ := strconv.Atoi(params["id"])

		if user.ID == id {
			json.NewEncoder(w).Encode(user)

			return
		}
	}

	json.NewEncoder(w).Encode(&User{})
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for index, value := range users {
		id, _ := strconv.Atoi(params["id"])

		if value.ID == id {
			var user User
			_ = json.NewDecoder(r.Body).Decode(&user)

			user.ID = users[index].ID

			users[index].Name = user.Name
			users[index].Surname = user.Surname
			users[index].Age = user.Age

			json.NewEncoder(w).Encode(user)

			return
		}
	}

	json.NewEncoder(w).Encode(users)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for index, value := range users {
		id, _ := strconv.Atoi(params["id"])

		if value.ID == id {
			users = append(users[:index], users[index+1:]...)

			break
		}
	}

	json.NewEncoder(w).Encode(users)
}
