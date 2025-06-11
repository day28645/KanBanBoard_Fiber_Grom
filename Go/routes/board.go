package routes

import (
	"errors"
	"fmt"
	"gofiber/database"
	"gofiber/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Board struct {
	gorm.Model
	User      User
	BoardName string `json:"board_name"`
	OwnerID   uint   `json:"owner_id"`
}

func createResponseBoard(board models.Board) Board {
	return Board{
		OwnerID:   board.User.ID,
		BoardName: board.BoardName,
	}
}

// POST
func CreateBoard(c *fiber.Ctx) error {
	var boardInput models.Board

	//parsing validation
	if err := c.BodyParser(&boardInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
			"data":    err.Error(),
		})
	}

	//input validation
	result := database.DB.
		Where("board_name = ?", boardInput.BoardName).
		Find(&boardInput)
	if result.RowsAffected > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Duplicated Board Name",
		})
	}

	if boardInput.BoardName == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid board name",
			"data":    nil,
		})
	}

	board := models.Board{
		BoardName: boardInput.BoardName,
		OwnerID:   boardInput.OwnerID,
	}

	database.DB.Create(&board)
	responseBoard := createResponseBoard(board)
	return c.Status(200).JSON(responseBoard)
}

// GET All Board
func GetBoards(c *fiber.Ctx) error {
	boards := []models.Board{}

	database.DB.Find(&boards)
	responseBoards := []Board{}

	for _, board := range boards {
		responseBoard := createResponseBoard(board)
		responseBoards = append(responseBoards, responseBoard)
	}
	return c.Status(200).JSON(responseBoards)
}

// query to find User in DB
func findBoard(id int, board *models.Board) error {
	database.DB.First(&board, "id=?", id)
	if board.ID == 0 {
		return errors.New("Board does not exist")
	}
	return nil
}

// GET by ID
func GetBoardByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var board models.Board

	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer")
	}

	if err := findBoard(id, &board); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	responseBoard := createResponseBoard(board)
	return c.Status(200).JSON(responseBoard)
}

// PUT
func UpdateBoard(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var board models.Board

	//check board id
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer")
	}

	//query to find Board
	if err := findBoard(id, &board); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateBoard struct {
		BoardName string `json:"board_name"`
	}

	var updateData UpdateBoard

	//parsing validation
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	//duplicated validation
	if updateData.BoardName != "" {
		var existingBoard models.Board
		result := database.DB.
			Where("board_name = ? AND id != ?", updateData.BoardName, board.ID).
			First(&existingBoard)

		if result.RowsAffected > 0 {
			return c.Status(400).JSON(fiber.Map{
				"error":   true,
				"message": "Duplicated board name",
			})
		}
		board.BoardName = updateData.BoardName
	}

	//null validation - if null, data is still the same
	if updateData.BoardName != "" {
		board.BoardName = updateData.BoardName
	}

	//update database
	database.DB.Save(&board)

	responseBoard := createResponseBoard(board)
	return c.Status(200).JSON(responseBoard)
}

// DELETE
func DeleteBoard(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	fmt.Println(id)
	var board models.Board

	//check board id
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer")
	}

	//query to find Board
	if err := findBoard(id, &board); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	//soft delete
	if err := database.DB.Delete(&board).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("Successfully Deleted Board")
}
