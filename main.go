package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

type Employee struct {
	ID     string `json:"id, omitempty" bson:"_id, omitempty"`
	Name   string `json:"name"`
	Salary string `json:"salary"`
	Age    string `json:"age"`
}

var mg MongoInstance

const dbName = "fiber-hrms"
const mongoURI = "mongodb://localhost:27017" + dbName

/**
 *
 * db connection func
 **/
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

/**
 *
 * get all employess records
 **/
func getEmployees(c *fiber.Ctx) error {

	query := bson.D{{}}
	cursor, err := mg.Db.Collection("Employees").Find(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var employees []Employee = make([]Employee, 0)

	if err := cursor.All(c.Context(), &employees); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(employees)

}

/**
 *
 *
 **/
func getEmployeeByID(c *fiber.Ctx) error {}

/**
 *
 *
 **/
func createNewEmployee(c *fiber.Ctx) error {}

/**
 *
 *
 **/
func updateEmployee(c *fiber.Ctx) error {}

/**
 *
 *
 **/
func deleteEmploye(c *fiber.Ctx) error {}

/**
 *
 *
 **/
func main() {

	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/employee", getEmployees)
	app.Get("/employee/:id", getEmployeeByID)
	app.Post("/employee", createNewEmployee)
	app.Put("/employee/:id", updateEmployee)
	app.Delete("/employee/:id", deleteEmploye)
}
