package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Database struct {
		Uri        string `yaml:"uri"`
		Host       string `yaml:"host"`
		Port       string `yaml:"port"`
		Db         string `yaml:"db"`
		Migrations string `yaml:"migrations"`
	} `yaml:"database"`
}

var config Config

func main() {
	acquireConfig(&config)

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s [up|down]", os.Args[0])
	}

	option := os.Args[1]

	log.SetFormatter(&log.JSONFormatter{})

	migrationsPath := "file://" + config.Database.Migrations

	databaseURI := strings.ReplaceAll(config.Database.Uri, "<host>", config.Database.Host)
	databaseURI = strings.ReplaceAll(databaseURI, "<port>", config.Database.Port)
	databaseURI += config.Database.Db

	clientOptions := options.Client().ApplyURI(databaseURI)

	migrator, err := migrate.New(migrationsPath, databaseURI)

	if err != nil {
		log.Fatalf("Migrator set up failed [%v]", err)
	}

	switch option {
	case "up":
		err = migrator.Up()

		if err != nil && err.Error() != "no change" {
			log.Fatalf("up failed [%v]", err)
		} else {
			log.Print("up done")
		}
	case "down":
		err = migrator.Down()

		if err != nil {
			log.Fatalf("down failed [%v]", err)
		} else {
			log.Print("down done")
		}
	default:
		log.Fatalf("Option invalid: %s", option)
	}

	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Fatalf("MongoDB client set up process failed [%v]", err)
	}

	dbContext, _ := context.WithTimeout(context.Background(), 7*time.Second) //TODO: context.CancelFunc

	err = client.Connect(dbContext)

	if err != nil {
		log.Fatalf("Unable to establish a connection to MongoDB [%v", err)
	}
	defer client.Disconnect(dbContext)

	databases, err := client.ListDatabaseNames(dbContext, bson.M{})

	if err != nil {
		log.Fatalf("Unable to read current databases [%v]", err)
	} else {
		log.Printf("Current databases: %v", databases)
	}
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
