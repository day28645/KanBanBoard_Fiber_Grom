package routes

import (
	"errors"
	"fmt"
	"gofiber/database"
	"gofiber/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BoardMember struct {
	gorm.Model
	User    User
	Board   Board
	Role    string `json:"role"`
	BoardID uint   `json:"board_id"`
	UserID  uint   `json:"user_id"`
}

func createResponseBoardMember(boardmember models.BoardMember) BoardMember {
	return BoardMember{
		BoardID: boardmember.Board.ID,
		UserID:  boardmember.UserID,
		Role:    boardmember.Role,
	}
}

// Post
func CreateBoardMember(c *fiber.Ctx) error {
	var boardMemberInput models.BoardMember

	//parsing validation
	if err := c.BodyParser(&boardMemberInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
			"data":    err.Error(),
		})
	}

	//input validation
	if boardMemberInput.Role == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid role",
			"data":    nil,
		})
	}

	//check if board exists
	var board models.Board
	if err := database.DB.First(&board, boardMemberInput.BoardID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "Board not found",
		})
	}

	//check if user exists
	var user models.User
	if err := database.DB.First(&user, boardMemberInput.UserID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "User not found",
		})
	}

	//validate role input
	validRoles := map[string]bool{
		"preparer": true,
		"reviewer": true,
		"viewer":   true,
	}
	if !validRoles[boardMemberInput.Role] {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid role. Must be one of: preparer, reviewer, viewer",
		})
	}

	//check assigning reviewer and ensure only one exists
	if boardMemberInput.Role == "reviewer" {
		var existingReviewer models.BoardMember
		err := database.DB.
			Where("board_id = ? AND role = ?", boardMemberInput.BoardID, "reviewer").
			First(&existingReviewer).Error

		if err == nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   true,
				"message": "This board already has a reviewer",
			})
		}
	}

	newBoard := models.BoardMember{
		BoardID: boardMemberInput.BoardID,
		UserID:  boardMemberInput.UserID,
		Role:    boardMemberInput.Role,
	}

	//insert database
	database.DB.Create(&newBoard)
	responseBoardMember := createResponseBoardMember(newBoard)
	return c.Status(200).JSON(responseBoardMember)
}

// GET All BoardMember
func GetBoardMembers(c *fiber.Ctx) error {
	boardMemberInput := []models.BoardMember{}

	//read database
	database.DB.Find(&boardMemberInput)
	responseBoardMembers := []BoardMember{}

	for _, boardmember := range boardMemberInput {
		responseBoardMember := createResponseBoardMember(boardmember)
		responseBoardMembers = append(responseBoardMembers, responseBoardMember)
	}
	return c.Status(200).JSON(responseBoardMembers)
}

// query to find BoardMember in DB
func findBoardMember(id int, boardmember *models.BoardMember) error {
	database.DB.First(&boardmember, "id=?", id)
	if boardmember.ID == 0 {
		return errors.New("Board Member does not exist")
	}
	return nil
}

// GET by ID
func GetBoardMemberByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var boardMemberInput models.BoardMember

	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer")
	}

	//read database
	if err := findBoardMember(id, &boardMemberInput); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	responseBoardMember := createResponseBoardMember(boardMemberInput)
	return c.Status(200).JSON(responseBoardMember)
}

// PUT
func UpdateBoardMember(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var boardMemberInput models.BoardMember

	//check board id
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer")
	}

	//query to find BoardMember
	if err := findBoardMember(id, &boardMemberInput); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateBoardMember struct {
		Role string `json:"role"`
	}

	var updateData UpdateBoardMember

	//parsing validation
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	//validate role input
	validRoles := map[string]bool{
		"preparer": true,
		"reviewer": true,
		"viewer":   true,
	}
	if !validRoles[updateData.Role] {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid role. Must be one of: preparer, reviewer, viewer",
		})
	}

	//check assigning reviewer and ensure only one exists
	if updateData.Role == "reviewer" {
		var count int64
		database.DB.Model(&models.BoardMember{}).
			Where("board_id = ? AND role = ? AND id != ?", boardMemberInput.BoardID, "reviewer", boardMemberInput.ID).
			Count(&count)

		if count > 0 {
			return c.Status(400).JSON(fiber.Map{
				"error":   true,
				"message": "Only one reviewer is allowed per board",
			})
		}
	}

	//null validation - if null, data is still the same
	if updateData.Role != "" {
		boardMemberInput.Role = updateData.Role
	}

	//update database
	database.DB.Save(&boardMemberInput)

	responseBoardMember := createResponseBoardMember(boardMemberInput)
	return c.Status(200).JSON(responseBoardMember)
}

// DELETE
func DeleteBoardMember(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	fmt.Println(id)
	var boardMemberInput models.BoardMember

	//check boardmember id
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer")
	}

	//query to find Board
	if err := findBoardMember(id, &boardMemberInput); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	//soft delete
	if err := database.DB.Delete(&boardMemberInput).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("Successfully Deleted BoardMember")
}
