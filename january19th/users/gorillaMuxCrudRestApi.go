package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"database"`
}

// User represents the user for this application
//
// swagger:model
type User struct {
	// the id for this user
	// required: true
	ID string `json:"id"`

	// the name for this user
	// required: true
	// min length: 3
	Name string `json:"name"`

	// the surname for this user
	// required: true
	// min length: 3
	Surname string `json:"surname"`

	// the age for this user
	// required: true
	// min: 1
	// max: 120
	Age int `json:"age"`
}

var config Config
var users []User

func main() {
	setupAndRunUsersServer()
}

func setupAndRunUsersServer() {
	acquireConfig(&config)
	serverHostPort := config.Server.Host + ":" + config.Server.Port

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
	log.Fatal(http.ListenAndServe(serverHostPort, router))
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

	databaseHostPort := config.Database.Host + ":" + config.Database.Port
	databaseURI := "mongodb://" + databaseHostPort + "/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+1.2.3"

	client, err := mongo.NewClient(options.Client().ApplyURI(databaseURI))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("MongoDB client ready to use")
	}
	ctx, _ := context.WithTimeout(context.Background(), 7*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("MongoDB client connected to %s\n", databaseHostPort)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("MongoDB client connection test succeeded")
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("MongoDB client acquired database names")
	}
	fmt.Println("These are the current databases:", databases)
}

func getId() string {
	uuid := uuid.NewString()

	return uuid
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /api/users createUser
	//
	// Creates a new user into the system that the user has access to
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: Body
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/User"
	// responses:
	//   '200':
	//     description: createUser response
	w.Header().Set("Content-Type", "application/json")

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	user.ID = getId()

	users = append(users, user)

	json.NewEncoder(w).Encode(user)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /api/users getUsers
	//
	// Returns all users from the system that the user has access to
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: getUsers response
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/User"
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /api/users/{id} getUser
	//
	// Returns a specific user from the system that the user has access to
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: id to filter by
	//   required: true
	//   type: integer
	// responses:
	//   '200':
	//     description: getUser response
	//     schema:
	//       "$ref": "#/definitions/User"
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
	// swagger:operation PUT /api/users/{id} updateUser
	//
	// Updates a specific user from the system that the user has access to
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: id to filter by
	//   required: true
	//   type: integer
	// responses:
	//   '200':
	//     description: updateUser response
	//     schema:
	//       "$ref": "#/definitions/User"
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
	// swagger:operation DELETE /api/users/{id} deleteUser
	//
	// Deletes a specific user from the system that the user has access to
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: id to filter by
	//   required: true
	//   type: integer
	// responses:
	//   '200':
	//     description: deleteUser response
	//     schema:
	//       "$ref": "#/definitions/User"
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
