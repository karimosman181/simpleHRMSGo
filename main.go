package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

type Employee struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Salary string `json:salary`
	Age    string `age`
}

var mg MongoInstance

const dbName = "fiber-hrms"
const mongoURI = "mongodb://localhost:27017" + dbName

func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	//define timeOut (so you donot block the entire program)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {
		return err
	}

	mg = MongoInstance{
		Client: client,
		Db:     db,
	}

	return nil
}

func main() {

	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/employee", getEmployees())
	app.Get("/employee/:id", getEmployeeByID())
	app.Post("/employee", createNewEmployee())
	app.Put("/employee/:id", updateEmployee())
	app.Delete("/employee/:id", deleteEmploye())
}
