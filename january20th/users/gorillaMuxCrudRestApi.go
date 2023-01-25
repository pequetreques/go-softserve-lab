package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Uri  string `yaml:"uri"`
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Db   string `yaml:"db"`
		Coll string `yaml:"coll"`
	} `yaml:"database"`
}

// User represents the user for this application
//
// swagger:model
type User struct {
	// the id for this user
	// required: true
	Id string `json:"id"`

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
var clientOptions *options.ClientOptions
var databaseURI string
var err error
var client *mongo.Client
var collection *mongo.Collection

func main() {
	acquireConfig(&config)
	setupUsersDatabase()
	setupAndRunUsersServer()
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

func getId() string {
	uuid := uuid.NewString()

	return uuid
}

func setupUsersDatabase() {
	databaseURI = strings.ReplaceAll(config.Database.Uri, "<host>", config.Database.Host)
	databaseURI = strings.ReplaceAll(databaseURI, "<port>", config.Database.Port)
	clientOptions = options.Client().ApplyURI(databaseURI)

	client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collection = client.Database(config.Database.Db).Collection(config.Database.Coll)

	documents := []interface{}{
		User{Id: getId(), Name: "Eren", Surname: "Jaeger", Age: 17},
		User{Id: getId(), Name: "Mikasa", Surname: "Ackerman", Age: 17},
		User{Id: getId(), Name: "Armin", Surname: "Arlert", Age: 16},
		User{Id: getId(), Name: "Levi", Surname: "Ackerman", Age: 27},
		User{Id: getId(), Name: "Annie", Surname: "Leonhart", Age: 18},
		User{Id: getId(), Name: "Christa", Surname: "Lenz", Age: 15},
		User{Id: getId(), Name: "Sasha", Surname: "Brouse", Age: 16},
		User{Id: getId(), Name: "Zoë", Surname: "Hange", Age: 29},
		User{Id: getId(), Name: "Gabi", Surname: "Braun", Age: 12},
	}

	collection.InsertMany(context.TODO(), documents)
}

func setupAndRunUsersServer() {
	serverHostPort := config.Server.Host + ":" + config.Server.Port

	log.Println("Initializing server...")
	router := mux.NewRouter()
	log.Println("...done")

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

	user.Id = getId()

	client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collection = client.Database(config.Database.Db).Collection(config.Database.Coll)

	_, err := collection.InsertOne(context.TODO(), user)

	if err != nil {
		json.NewEncoder(w).Encode(&User{})
	} else {
		json.NewEncoder(w).Encode(user)
	}
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

	client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collection = client.Database(config.Database.Db).Collection(config.Database.Coll)

	filter := bson.D{}

	cursor, err := collection.Find(context.TODO(), filter)

	if err != nil {
		json.NewEncoder(w).Encode(&User{})
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		json.NewEncoder(w).Encode(&User{})
	}

	cursor.Close(context.TODO())
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

	var user User

	params := mux.Vars(r)

	client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collection = client.Database(config.Database.Db).Collection(config.Database.Coll)

	filter := bson.D{{Key: "id", Value: params["id"]}}

	err = collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		json.NewEncoder(w).Encode(&User{})
	} else {
		json.NewEncoder(w).Encode(user)
	}
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

	var user User

	params := mux.Vars(r)
	id := params["id"]

	_ = json.NewDecoder(r.Body).Decode(&user)

	client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collection = client.Database(config.Database.Db).Collection(config.Database.Coll)

	filter := bson.D{{Key: "id", Value: id}}

	replacement := User{Id: id, Name: user.Name, Surname: user.Surname, Age: user.Age}

	_, err = collection.ReplaceOne(context.TODO(), filter, replacement)

	if err != nil {
		json.NewEncoder(w).Encode(&User{})
	} else {
		user.Id = id
		json.NewEncoder(w).Encode(user)
	}
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

	client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collection = client.Database(config.Database.Db).Collection(config.Database.Coll)

	filter := bson.D{{Key: "id", Value: params["id"]}}

	_, err = collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		panic(err)
	}

	filter = bson.D{}

	cursor, err := collection.Find(context.TODO(), filter)

	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(users)
}
