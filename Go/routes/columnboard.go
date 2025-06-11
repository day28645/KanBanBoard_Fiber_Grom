package routes

import (
	"errors"
	"fmt"
	"gofiber/database"
	"gofiber/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ColumnBoard struct {
	gorm.Model
	Board      Board
	ColumnName string `json:"column_name"`
	BoardID    uint   `json:"board_id"`
}

func createResponseColumnBoard(columnboard models.ColumnBoard) ColumnBoard {
	return ColumnBoard{
		BoardID:    columnboard.Board.ID,
		ColumnName: columnboard.ColumnName,
	}
}

// POST
func CreateColumnBoard(c *fiber.Ctx) error {
	var columnboardInput models.ColumnBoard

	//parsing validation
	if err := c.BodyParser(&columnboardInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
			"data":    err.Error(),
		})
	}

	//check if board exists
	var board models.Board
	if err := database.DB.First(&board, columnboardInput.BoardID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "Board not found",
		})
	}

	//validate column input
	validColumns := map[string]bool{
		"To Do":    true,
		"Doing":    true,
		"Done":     true,
		"Accepted": true,
	}
	if !validColumns[columnboardInput.ColumnName] {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid column. Must be one of: To Do, Doing, Done, Accepted",
		})
	}

	//input validation
	result := database.DB.
		Where("board_id = ? AND column_name = ?", columnboardInput.BoardID, columnboardInput.ColumnName).
		Find(&columnboardInput)
	if result.RowsAffected > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "This column already exists in the board",
		})
	}

	newColumn := models.ColumnBoard{
		BoardID:    columnboardInput.BoardID,
		ColumnName: columnboardInput.ColumnName,
	}
	//insert database
	database.DB.Create(&newColumn)
	responseColumnBoard := createResponseColumnBoard(newColumn)
	return c.Status(200).JSON(responseColumnBoard)
}

// GET All BoardMember
func GetColumnBoards(c *fiber.Ctx) error {
	columnboardInput := []models.ColumnBoard{}

	//read database
	database.DB.Find(&columnboardInput)
	responseColumnBoards := []ColumnBoard{}

	for _, columnboard := range columnboardInput {
		responseColumnBoard := createResponseColumnBoard(columnboard)
		responseColumnBoards = append(responseColumnBoards, responseColumnBoard)
	}
	return c.Status(200).JSON(responseColumnBoards)
}

// query to find ColumnBoard in DB
func findColumnBoard(id int, columnboard *models.ColumnBoard) error {
	database.DB.First(&columnboard, "id=?", id)
	if columnboard.ID == 0 {
		return errors.New("Column Board does not exist")
	}
	return nil
}

// GET by ID
func GetColumnBoardByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var columnboard models.ColumnBoard

	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer")
	}

	//read database
	if err := findColumnBoard(id, &columnboard); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	responseColumnBoard := createResponseColumnBoard(columnboard)
	return c.Status(200).JSON(responseColumnBoard)
}

// PUT
func UpdateColumnBoard(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var columnboardInput models.ColumnBoard

	//check column id
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer")
	}

	//query to find ColumnBoard
	if err := findColumnBoard(id, &columnboardInput); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateColumnBoard struct {
		ColumnName string `json:"column_name"`
	}

	var updateData UpdateColumnBoard

	//parsing validation
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	//validate role input
	validColumns := map[string]bool{
		"To Do":    true,
		"Doing":    true,
		"Done":     true,
		"Accepted": true,
	}
	if !validColumns[updateData.ColumnName] {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid column. Must be one of: To Do, Doing, Done, Accepted",
		})
	}

	//input validation
	var existingBoard models.ColumnBoard
	result := database.DB.
		Where("board_id = ? AND column_name = ? AND id != ?", columnboardInput.BoardID, updateData.ColumnName, columnboardInput.ID).
		First(&existingBoard)

	if result.RowsAffected > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "This column already exists in the board",
		})
	}

	//null validation - if null, data is still the same
	if updateData.ColumnName != "" {
		columnboardInput.ColumnName = updateData.ColumnName
	}

	//update database
	database.DB.Save(&columnboardInput)

	responseColumnBoard := createResponseColumnBoard(columnboardInput)
	return c.Status(200).JSON(responseColumnBoard)
}

// DELETE
func DeleteColumnBoard(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	fmt.Println(id)
	var columnboardInput models.ColumnBoard

	//check boardmember id
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer")
	}

	//query to find Board
	if err := findColumnBoard(id, &columnboardInput); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	//soft delete
	if err := database.DB.Delete(&columnboardInput).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("Successfully Deleted Column Board")
}
