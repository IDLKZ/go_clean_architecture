package api_routing

import (
	"clean_architecture_fiber/app/route/handler"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoleRoutes(app *fiber.App, roleHandler *handler.RoleHandler) {
	api := app.Group("/api/v1")
	roles := api.Group("/roles")
	roles.Get("/:value", roleHandler.GetByValue)
}
