package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
}

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     int    `json:"age"`
}

var config Config
var users []User

func main() {
	setupAndRunUsersServer()
}

func setupAndRunUsersServer() {
	acquireConfig(&config)
	hostPort := config.Server.Host + ":" + config.Server.Port

	log.Println("Initializing server...")
	router := mux.NewRouter()
	log.Println("...done")

	setupDatabase()

	log.Println("Configuring server...")
	router.HandleFunc("/api/users", createUser).Methods("POST")
	router.HandleFunc("/api/users", getUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/api/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/api/users/{id}", deleteUser).Methods("DELETE")
	log.Println("...done")

	log.Println("Running server...")
	log.Fatal(http.ListenAndServe(hostPort, router))
	log.Println("...done")
}

func acquireConfig(c *Config) {
	yamlConfigFile, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal("Could not acquire server configuration")
	}
	defer yamlConfigFile.Close()

	decoder := yaml.NewDecoder(yamlConfigFile)
	err = decoder.Decode(&c)
	if err != nil {
		log.Fatal("Could not acquire server configuration")
	}
}

func setupDatabase() {
	users = append(users, User{ID: getId(), Name: "Eren", Surname: "Jaeger", Age: 17})
	users = append(users, User{ID: getId(), Name: "Mikasa", Surname: "Ackerman", Age: 17})
	users = append(users, User{ID: getId(), Name: "Armin", Surname: "Arlert", Age: 16})
	users = append(users, User{ID: getId(), Name: "Levi", Surname: "Ackerman", Age: 27})
	users = append(users, User{ID: getId(), Name: "Annie", Surname: "Leonhart", Age: 18})
	users = append(users, User{ID: getId(), Name: "Christa", Surname: "Lenz", Age: 15})
	users = append(users, User{ID: getId(), Name: "Sasha", Surname: "Brouse", Age: 16})
	users = append(users, User{ID: getId(), Name: "Zoë", Surname: "Hange", Age: 29})
}

func getId() string {
	uuid := uuid.NewString()

	return uuid
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	user.ID = getId()

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
		if user.ID == params["id"] {
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
		if value.ID == params["id"] {
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
		if value.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)

			break
		}
	}

	json.NewEncoder(w).Encode(users)
}
