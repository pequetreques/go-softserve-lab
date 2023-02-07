package main

import (
	"context"
	"log"
	"net"
	"os"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
	pb "softserve.com/january25th/generated/proto"

	"google.golang.org/grpc"
)

type usersApiServer struct {
	pb.UnimplementedUsersApiServer
}

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

type User struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     int    `json:"age"`
}

var config Config
var users []User
var clientOptions *options.ClientOptions
var databaseURI string
var err error
var client *mongo.Client
var collection *mongo.Collection
var listener net.Listener
var server *grpc.Server

func main() {
	acquireConfig(&config)
	setupUsersDatabase()
	setupAndRunUsersServer()
}

func acquireConfig(c *Config) {
	yamlConfigFile, err := os.Open("config/config.yaml")
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

	listener, err = net.Listen("tcp", serverHostPort)

	if err != nil {
		log.Fatalln(err)
	}

	server = grpc.NewServer()
	pb.RegisterUsersApiServer(server, &usersApiServer{})

	err = server.Serve(listener)

	if err != nil {
		log.Fatalln(err)
	}
}

func (s *usersApiServer) CreateUser(ctx context.Context, request *pb.CreateRequest) (*pb.UserResponse, error) {
	var user User
	var userResponse pb.UserResponse

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

	user.Id = getId()
	user.Name = request.GetName()
	user.Surname = request.GetSurname()
	user.Age = int(request.GetAge())

	_, err := collection.InsertOne(context.TODO(), user)

	if err != nil {
		return &pb.UserResponse{}, nil
	} else {
		userResponse.Id = user.Id
		userResponse.Name = user.Name
		userResponse.Surname = user.Surname
		userResponse.Age = int32(user.Age)

		return &userResponse, nil
	}
}

func (s *usersApiServer) GetUsers(request *pb.EmptyRequest, response pb.UsersApi_GetUsersServer) error {
	var userResponse pb.UserResponse

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
		return nil
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil
	}

	cursor.Close(context.TODO())

	for _, user := range users {
		userResponse.Id = user.Id
		userResponse.Name = user.Name
		userResponse.Surname = user.Surname
		userResponse.Age = int32(user.Age)

		if err = response.Send(&userResponse); err != nil {
			return err
		}
	}

	return nil
}

func (s *usersApiServer) GetUser(ctx context.Context, request *pb.IdRequest) (*pb.UserResponse, error) {
	var user User
	var userResponse pb.UserResponse

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

	filter := bson.D{{Key: "id", Value: request.GetId()}}

	err = collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return &pb.UserResponse{}, nil
	} else {
		userResponse.Id = user.Id
		userResponse.Name = user.Name
		userResponse.Surname = user.Surname
		userResponse.Age = int32(user.Age)

		return &userResponse, nil
	}
}

func (s *usersApiServer) UpdateUser(ctx context.Context, request *pb.IdBodyRequest) (*pb.UserResponse, error) {
	var userResponse pb.UserResponse

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

	filter := bson.D{{Key: "id", Value: request.GetId()}}

	replacement := User{Id: request.GetId(), Name: request.GetName(), Surname: request.GetSurname(), Age: int(request.GetAge())}

	_, err = collection.ReplaceOne(context.TODO(), filter, replacement)

	if err != nil {
		return &pb.UserResponse{}, nil
	} else {
		userResponse.Id = request.GetId()
		userResponse.Name = request.GetName()
		userResponse.Surname = request.GetSurname()
		userResponse.Age = request.GetAge()

		return &userResponse, nil
	}
}

func (s *usersApiServer) DeleteUser(ctx context.Context, request *pb.IdRequest) (*pb.UserResponse, error) {
	var user User
	var userResponse pb.UserResponse

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

	filter := bson.D{{Key: "id", Value: request.GetId()}}

	err = collection.FindOne(context.TODO(), filter).Decode(&user)

	userResponse.Id = user.Id
	userResponse.Name = user.Name
	userResponse.Surname = user.Surname
	userResponse.Age = int32(user.Age)

	filter = bson.D{{Key: "id", Value: request.GetId()}}

	_, err = collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		return &pb.UserResponse{}, nil
	}

	return &userResponse, nil
}
