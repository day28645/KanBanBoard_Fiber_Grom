package routes

import (
	"errors"
	"fmt"
	"gofiber/database"
	"gofiber/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TaskAssignee struct {
	gorm.Model
	Task             Task
	Assignee         User
	AssignedBy       User
	TaskID           uint `json:"task_id"`
	UserID           uint `json:"user_id"`
	AssignedByUserID uint `json:"assigned_by_user_id"`
}

func createResponseTaskAssignee(taskAssignee models.TaskAssignee) TaskAssignee {
	return TaskAssignee{
		TaskID:           taskAssignee.TaskID,
		UserID:           taskAssignee.UserID,
		AssignedByUserID: taskAssignee.AssignedByUserID,
	}
}

// POST
func CreateTaskAssignee(c *fiber.Ctx) error {
	var taskAssigneeInput models.TaskAssignee

	//parsing validation
	if err := c.BodyParser(&taskAssigneeInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
			"data":    err.Error(),
		})
	}

	//input validation
	if taskAssigneeInput.TaskID == 0 && taskAssigneeInput.UserID == 0 && taskAssigneeInput.AssignedByUserID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "All field is required",
			"data":    nil,
		})
	} else if taskAssigneeInput.TaskID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid task id",
			"data":    nil,
		})
	} else if taskAssigneeInput.UserID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid user id",
			"data":    nil,
		})
	} else if taskAssigneeInput.AssignedByUserID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid assgined by user id",
			"data":    nil,
		})
	}

	// check assigneeUser is exsist
	var assigneeUser models.User
	if err := database.DB.First(&assigneeUser, taskAssigneeInput.UserID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "Assignee user not found",
		})
	}

	// check assignedByUser is exsist
	var assignedByUser models.User
	if err := database.DB.First(&assignedByUser, taskAssigneeInput.AssignedByUserID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "Assigned-by user not found",
		})
	}

	//check title and ensure only one exists
	var count int64
	database.DB.Model(&models.TaskAssignee{}).
		Where("task_id = ? AND assigned_by_user_id = ?", taskAssigneeInput.TaskID, taskAssigneeInput.AssignedByUserID).
		Count(&count)

	if count > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "This task is already assigned by this user",
		})
	}

	newTaskAssignee := models.TaskAssignee{
		TaskID:           taskAssigneeInput.TaskID,
		UserID:           taskAssigneeInput.UserID,
		AssignedByUserID: taskAssigneeInput.AssignedByUserID,
	}

	//insert database
	database.DB.Create(&newTaskAssignee)
	responseTaskAssignee := createResponseTaskAssignee(newTaskAssignee)
	return c.Status(200).JSON(responseTaskAssignee)
}

// GET All BoardMember
func GetTaskAssignees(c *fiber.Ctx) error {
	taskAssignees := []models.TaskAssignee{}

	database.DB.Find(&taskAssignees)
	responseTaskAssignees := []TaskAssignee{}

	for _, taskAssignee := range taskAssignees {
		response := createResponseTaskAssignee(taskAssignee)
		responseTaskAssignees = append(responseTaskAssignees, response)
	}
	return c.Status(200).JSON(responseTaskAssignees)
}

func findTaskAssignee(id int, taskAssignee *models.TaskAssignee) error {
	database.DB.First(&taskAssignee, "id=?", id)
	if taskAssignee.ID == 0 {
		return errors.New("TaskAssignee does not exist")
	}
	return nil
}

// GET by ID
func GetTaskAssigneeByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var taskAssignee models.TaskAssignee

	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer")
	}

	if err := findTaskAssignee(id, &taskAssignee); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	responseTaskAssignee := createResponseTaskAssignee(taskAssignee)
	return c.Status(200).JSON(responseTaskAssignee)
}

// DELETE
func DeleteTaskAssignee(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	fmt.Println(id)
	var taskAssignee models.TaskAssignee

	//check board id
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer")
	}

	//query to find Board
	if err := findTaskAssignee(id, &taskAssignee); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	//soft delete
	if err := database.DB.Delete(&taskAssignee).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("Successfully Deleted TasktaskAssignee")
}
