package route

import (
	"clean_architecture_fiber/app/route/api_routing"
	"clean_architecture_fiber/app/route/handler"
	"github.com/gofiber/fiber/v2"
)

// setupRoutes настраивает маршруты API
func SetupRoutes(app *fiber.App, roleHandler *handler.RoleHandler) {
	api_routing.RegisterRoleRoutes(app, roleHandler)
}
