package di

import (
	"medassist/config"
	"medassist/internal/auth"
	"medassist/internal/nurse"
	"medassist/internal/repository"
	"medassist/internal/user"
	"medassist/internal/admin"
)

type Container struct {
	AuthHandler  *auth.AuthHandler
	UserHandler  *user.UserHandler
	NurseHandler *nurse.NurseHandler
	AdminHandler *admin.AdminHandler
}

func NewContainer() *Container {
	// Inicializa o banco de dados
	db := config.GetMongoDB()
	// Construtores: repository → service → handler
	userRepository := repository.NewUserRepository(db)
	nurseRepository := repository.NewNurseRepository(db)
	visitRepository := repository.NewVisitRepository(db)
	authService := auth.NewAuthService(userRepository, nurseRepository)
	adminService := admin.NewAdminService(userRepository, nurseRepository)
	userService := user.NewUserService(userRepository, nurseRepository, visitRepository)
	nurseService := nurse.NewNurseService(nurseRepository, visitRepository)

	authHandler := auth.NewAuthHandler(authService)
	adminHandler := admin.NewAdminHandler(adminService)
	userHandler := user.NewUserHandler(userService)
	nurseHandler := nurse.NewNurseHandler(nurseService)

	return &Container{
		AuthHandler:  authHandler,
		AdminHandler: adminHandler,
		UserHandler:  userHandler,
		NurseHandler: nurseHandler,
	}
}
