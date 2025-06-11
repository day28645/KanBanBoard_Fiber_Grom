package routes

import (
	"gofiber/database"
	"gofiber/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Login struct {
	gorm.Model
	User   User
	UserID uint `json:"user_id"`
}

func CreateResponseLogin(login models.Login) Login {
	return Login{
		UserID: login.UserID,
	}
}

// GET
func CreateLogin(c *fiber.Ctx) error {
	var userInput models.User

	//parsing validation
	if err := c.BodyParser(&userInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
			"data":    err.Error(),
		})
	}

	//input validation
	result := database.DB.
		Where("username = ? AND password = ?", userInput.Username, userInput.Password).
		First(&userInput)
	if result.RowsAffected == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Incorrect username or password",
		})
	}

	if userInput.Username == "" && userInput.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid username or password",
			"data":    nil,
		})
	} else if userInput.Username == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Username is required",
			"data":    nil,
		})
	} else if userInput.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Password is required",
			"data":    nil,
		})
	}

	login := models.Login{
		UserID: userInput.ID,
	}

	//insert database
	database.DB.Create(&login)
	responseLogin := CreateResponseLogin(login)
	return c.Status(200).JSON(responseLogin)
}
