package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
func getEmployeeByID(c *fiber.Ctx) error {
	//get id from params
	id := c.Params("id")

	query := bson.D{{Key: "_id", Value: id}}
	cursor := mg.Db.Collection("employees").FindOne(c.Context(), query)

	createdEmployee := &Employee{}

	cursor.Decode(createdEmployee)

	return c.Status(201).JSON(createdEmployee)
}

/**
 *
 * create a new employee record
 **/
func createNewEmployee(c *fiber.Ctx) error {
	collection := mg.Db.Collection("employees")

	employee := new(Employee)

	if err := c.BodyParser(employee); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	//to force mongo db to create id
	employee.ID = ""

	insertionResult, err := collection.InsertOne(c.Context(), employee)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	//to make sure the record is inserted , retrive the record and send it give as a response
	ID := insertionResult.InsertedID

	query := bson.D{{Key: "_id", Value: ID}}
	cursor := collection.FindOne(c.Context(), query)

	createdEmployee := &Employee{}

	cursor.Decode(createdEmployee)

	return c.Status(201).JSON(createdEmployee)
}

/**
 *
 * update employee
 **/
func updateEmployee(c *fiber.Ctx) error {
	//get id from params
	id := c.Params("id")

	employeeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	employee := new(Employee)

	if err := c.BodyParser(employee); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	query := bson.D{{Key: "_id", Value: employeeID}}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "name", Value: employee.Name},
				{Key: "salary", Value: employee.Salary},
				{Key: "age", Value: employee.Age},
			},
		},
	}

	err = mg.Db.Collection("employees").FindOneAndUpdate(c.Context(), query, update).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.SendStatus(400)
		}

		return c.SendStatus(500)
	}

	employee.ID = id

	return c.Status(200).JSON(employee)
}

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
