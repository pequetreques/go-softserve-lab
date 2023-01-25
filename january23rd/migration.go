package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/eminetto/clean-architecture-go/migrations"
	migrate "github.com/eminetto/mongo-migrate"
	"github.com/globalsign/mgo"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Uri        string `yaml:"uri"`
		Host       string `yaml:"host"`
		Port       string `yaml:"port"`
		Db         string `yaml:"db"`
		Coll       string `yaml:"coll"`
		Migrations string `yaml:"migrations"`
		Template   string `yaml:"template"`
	} `yaml:"database"`
}

var config Config

func main() {
	acquireConfig(&config)

	if len(os.Args) == 1 {
		log.Fatal("Missing options: up or down")
	}

	option := os.Args[1]

	databaseURI := strings.ReplaceAll(config.Database.Uri, "<host>", config.Database.Host)
	databaseURI = strings.ReplaceAll(databaseURI, "<port>", config.Database.Port)
	fmt.Println(databaseURI)

	session, err := mgo.Dial(databaseURI)

	if err != nil {
		log.Fatal(err.Error())
	}
	defer session.Close()

	db := session.DB(config.Database.Db)
	migrate.SetDatabase(db)
	migrate.SetMigrationsCollection(config.Database.Migrations)
	migrate.SetLogger(log.New(os.Stdout, "INFO: ", 0))

	switch option {
	case "new":
		if len(os.Args) != 3 {
			log.Fatal("Should be: new description-of-migration")
		}

		template := fmt.Sprintf("./%s/%s.go", config.Database.Migrations, config.Database.Template)
		migration := fmt.Sprintf("./%s/%s_%s.go", config.Database.Migrations, time.Now().Format("20230123001545"), os.Args[2])

		from, err := os.Open(template)

		if err != nil {
			log.Fatal("Template file is missing!")
		}
		defer from.Close()

		to, err := os.OpenFile(migration, os.O_RDWR|os.O_CREATE, 0666)

		if err != nil {
			log.Fatal(err.Error())
		}
		defer to.Close()

		_, err = io.Copy(to, from)

		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("New migration created: %s\n", migration)
	case "up":
		err = migrate.Up(migrate.AllAvailable)
	case "down":
		err = migrate.Down(migrate.AllAvailable)
	}

	if err != nil {
		log.Fatal(err.Error())
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
