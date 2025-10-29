package dependecy_injection

import (
	"clean_architecture_fiber/app/route/handler"
	"clean_architecture_fiber/domain/repositories"
	"clean_architecture_fiber/domain/use_case/role_use_case"
	"go.uber.org/fx"
)

// RoleModule — независимый DI-модуль для домена "Role"
var RoleModule = fx.Options(
	fx.Provide(
		repositories.NewRoleRepository,
		role_use_case.NewGetRoleByValueUseCase,
		handler.NewRoleHandler,
	),
)
