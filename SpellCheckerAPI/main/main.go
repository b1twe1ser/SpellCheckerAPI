package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(requestid.New())

	app.Get("/word/correct", getWord)

	app.Listen(":80")

}
