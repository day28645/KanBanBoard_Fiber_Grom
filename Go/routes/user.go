package routes

import (
	"errors"
	"fmt"
	"gofiber/database"
	"gofiber/models"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	IDCard   string `json:"id_card"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func CreateResponseUser(user models.User) User {
	return User{
		IDCard:   user.IDCard,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}
}

// Post
func CreateUser(c *fiber.Ctx) error {
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
	if userInput.IDCard == "" && userInput.Username == "" && userInput.Password == "" && userInput.Email == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "All field is required",
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
	} else if userInput.Email == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Email is required",
			"data":    nil,
		})
	}

	//duplicated validation
	result := database.DB.
		Where("id_card = ?", userInput.IDCard).
		First(&userInput)

	if result.RowsAffected > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Duplicated ID Card",
		})
	}

	//insert database
	database.DB.Create(&userInput)
	responseUser := CreateResponseUser(userInput)
	return c.Status(200).JSON(responseUser)
}

// GET All User
func GetUsers(c *fiber.Ctx) error {
	users := []models.User{}

	//read database
	database.DB.Find(&users)
	responseUsers := []User{}

	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}
	return c.Status(200).JSON(responseUsers)
}

// query to find User in DB
func findUser(id int, user *models.User) error {
	database.DB.First(&user, "id=?", id)
	if user.ID == 0 {
		return errors.New("User does not exist")
	}
	return nil
}

// GET by ID
func GetUserByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer")
	}

	//read database
	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)
}

// PUT
func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var user models.User

	//check user id
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer")
	}

	//query to find User
	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateUser struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	var updateData UpdateUser

	//parsing validation
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	//null validation - if null, data is still the same
	if updateData.Username != "" {
		user.Username = updateData.Username
	}
	if updateData.Password != "" {
		user.Password = updateData.Password
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}

	//update database
	database.DB.Save(&user)

	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)
}

// DELETE
func DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	fmt.Println(id)
	var user models.User

	//check user id
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer")
	}

	//query to find User
	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	//soft delete
	if err := database.DB.Delete(&user).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("Successfully Deleted User")
}
