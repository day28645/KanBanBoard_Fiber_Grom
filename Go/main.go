package main

import (
	"gofiber/database"
	"gofiber/routes"

	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	//users endpoints
	app.Post("/api/users", routes.CreateUser)
	app.Get("/api/users", routes.GetUsers)
	app.Get("/api/users/:id", routes.GetUserByID)
	app.Put("/api/users/:id", routes.UpdateUser)
	app.Delete("/api/users/:id", routes.DeleteUser)

	//login endpoints
	app.Post("/api/login", routes.CreateLogin)

	//boards endpoints
	app.Post("/api/boards", routes.CreateBoard)
	app.Get("/api/boards", routes.GetBoards)
	app.Get("/api/boards/:id", routes.GetBoardByID)
	app.Put("/api/boards/:id", routes.UpdateBoard)
	app.Delete("/api/boards/:id", routes.DeleteBoard)

	//boardmembers endpoints
	app.Post("/api/boardmembers", routes.CreateBoardMember)
	app.Get("/api/boardmembers", routes.GetBoardMembers)
	app.Get("/api/boardmembers/:id", routes.GetBoardMemberByID)
	app.Put("/api/boardmembers/:id", routes.UpdateBoardMember)
	app.Delete("/api/boardmembers/:id", routes.DeleteBoardMember)

	//columnboards endpoints
	app.Post("/api/columnboards", routes.CreateColumnBoard)
	app.Get("/api/columnboards", routes.GetColumnBoards)
	app.Get("/api/columnboards/:id", routes.GetColumnBoardByID)
	app.Put("/api/columnboards/:id", routes.UpdateColumnBoard)
	app.Delete("/api/columnboards/:id", routes.DeleteColumnBoard)

	//tasks endpoints
	app.Post("/api/tasks", routes.CreateTask)
	app.Get("/api/tasks", routes.GetTasks)
	app.Get("/api/tasks/:id", routes.GetTaskByID)
	app.Put("/api/tasks/:id", routes.UpdateTask)
	app.Delete("/api/tasks/:id", routes.DeleteTask)

	//taskassignees endpoints
	app.Post("/api/taskassignees", routes.CreateTaskAssignee)
	app.Get("/api/taskassignees", routes.GetTaskAssignees)
	app.Get("/api/taskassignees/:id", routes.GetTaskAssigneeByID)
	app.Delete("/api/taskassignees/:id", routes.DeleteTaskAssignee)
}

func main() {
	database.ConnectDB()

	app := fiber.New()
	setupRoutes(app)

	app.Listen(":8000")
}
