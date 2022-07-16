package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/mukul26s/Rest-API-using-Go/models"
	"github.com/mukul26s/Rest-API-using-Go/storage"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	//Repository
	DB *gorm.DB
}

//Repository Methods :

//createbook
func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	//converting context(request) into book interface
	//& not applied before - so error occured
	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"}, //sending response
		)
		return err
	}

	err = r.DB.Create(&book).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "couldn't create book"},
		)
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "book added",
		})
	return nil
}

//GetAllBooks
func (r *Repository) GetAllBooks(context *fiber.Ctx) error {
	bookmodels := &[]models.Books{}

	err := r.DB.Find(&bookmodels).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "couldn't get the books "})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "book fetched successfully",
			"data":    bookmodels,
		})
	return nil

}

//DeleteBook
func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id not given",
		})
		return nil
	}
	err := r.DB.Delete(bookModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "delete successful",
	})
	return nil
}

//GetBooksByID
func (r *Repository) GetBooksByID(context *fiber.Ctx) error {
	id := context.Params("id")
	//& not there before
	bookModel := &models.Books{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	fmt.Println("id is ", id)
	err := r.DB.Where("id = ?", id).First(bookModel).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not get books by id",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "id fetched successfully",
		"data":    bookModel,
	})
	return nil
}

//setting up routes
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/createbook", r.CreateBook)
	api.Delete("deletebook/:id", r.DeleteBook)
	api.Get("/getbook/:id", r.GetBooksByID)
	api.Get("/getbooks", r.GetAllBooks)

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBname:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("DB not loaded")
	}

	err = models.MigrateBooks(db)

	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
