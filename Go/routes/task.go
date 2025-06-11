package routes

import (
	"errors"
	"fmt"
	"gofiber/database"
	"gofiber/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ColumnBoard    ColumnBoard
	User           User
	Title          string    `json:"title"`
	DueDate        time.Time `json:"due_date"`
	ColumnBoardID  uint      `json:"column_board_id"`
	CreateByUserID uint      `json:"create_by_user_id"`
}

func createResponseTask(task models.Task) Task {
	return Task{
		ColumnBoardID:  task.ColumnBoardID,
		CreateByUserID: task.CreateByUserID,
		Title:          task.Title,
		DueDate:        task.DueDate,
	}
}

// POST
func CreateTask(c *fiber.Ctx) error {
	var taskInput models.Task

	//parsing validation
	if err := c.BodyParser(&taskInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
			"data":    err.Error(),
		})
	}

	//input validation
	if taskInput.Title == "" && taskInput.DueDate.IsZero() && taskInput.ColumnBoardID == 0 && taskInput.CreateByUserID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "All field is required",
			"data":    nil,
		})
	} else if taskInput.Title == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid title",
			"data":    nil,
		})
	} else if taskInput.DueDate.IsZero() {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid due date",
			"data":    nil,
		})
	} else if taskInput.ColumnBoardID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid column board",
			"data":    nil,
		})
	} else if taskInput.CreateByUserID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid Create User ID",
			"data":    nil,
		})
	}

	//check if user exists
	var user models.User
	if err := database.DB.First(&user, taskInput.CreateByUserID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "User not found",
		})
	}

	//validate title input
	validTitle := map[string]bool{
		"New":         true,
		"In Progress": true,
		"Completed":   true,
	}
	if !validTitle[taskInput.Title] {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid role. Must be one of: New, In Progress, Completed",
		})
	}

	//check title and ensure only one exists
	result := database.DB.
		Where("column_board_id = ? AND title = ?", taskInput.ColumnBoardID, taskInput.Title).
		Find(&taskInput)
	if result.RowsAffected > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "This task already exists in the column board",
		})
	}

	newTask := models.Task{
		ColumnBoardID:  taskInput.ColumnBoardID,
		CreateByUserID: taskInput.CreateByUserID,
		Title:          taskInput.Title,
		DueDate:        taskInput.DueDate,
	}

	//insert database
	database.DB.Create(&newTask)
	responseTask := createResponseTask(newTask)
	return c.Status(200).JSON(responseTask)
}

// GET All BoardMember
func GetTasks(c *fiber.Ctx) error {
	tasks := []models.Task{}

	database.DB.Find(&tasks)
	responseTasks := []Task{}

	for _, task := range tasks {
		responseTask := createResponseTask(task)
		responseTasks = append(responseTasks, responseTask)
	}
	return c.Status(200).JSON(responseTasks)
}

// query to find Task in DB
func findTask(id int, task *models.Task) error {
	database.DB.First(&task, "id=?", id)
	if task.ID == 0 {
		return errors.New("Task does not exist")
	}
	return nil
}

// GET by ID
func GetTaskByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var task models.Task

	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer")
	}

	if err := findTask(id, &task); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	responseTask := createResponseTask(task)
	return c.Status(200).JSON(responseTask)
}

// PUT
func UpdateTask(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var taskInput models.Task

	//check task id
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer")
	}

	//query to find Task
	if err := findTask(id, &taskInput); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateTask struct {
		Title string `json:"title"`
	}

	var updateData UpdateTask

	//parsing validation
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	//check if task exists
	if err := database.DB.First(&taskInput, taskInput.ID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "User not found",
		})
	}

	//validate title input
	validTitle := map[string]bool{
		"New":         true,
		"In Progress": true,
		"Completed":   true,
	}
	if updateData.Title != "" && !validTitle[updateData.Title] {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid role. Must be one of: New, In Progress, Completed",
		})
	}

	//duplicated validation
	if updateData.Title != "" {
		var existingTask models.Task
		result := database.DB.
			Where("column_board_id = ? AND title = ?", taskInput.ColumnBoardID, updateData.Title).
			First(&existingTask)

		if result.RowsAffected > 0 {
			return c.Status(400).JSON(fiber.Map{
				"error":   true,
				"message": "Duplicated task",
			})
		}
		taskInput.Title = updateData.Title
	}

	//null validation - if null, data is still the same
	if updateData.Title != "" {
		taskInput.Title = updateData.Title
	}

	//update database
	database.DB.Save(&taskInput)

	responseTask := createResponseTask(taskInput)
	return c.Status(200).JSON(responseTask)
}

// DELETE
func DeleteTask(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	fmt.Println(id)
	var task models.Task

	//check board id
	if err != nil {
		return c.Status(400).JSON("Please ensure that id is an integer")
	}

	//query to find Board
	if err := findTask(id, &task); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	//soft delete
	if err := database.DB.Delete(&task).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("Successfully Deleted Task")
}
