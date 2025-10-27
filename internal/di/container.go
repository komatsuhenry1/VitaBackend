package di

import (
	"medassist/config"
	"medassist/internal/admin"
	"medassist/internal/auth"
	"medassist/internal/chat"
	"medassist/internal/nurse"
	"medassist/internal/repository"
	"medassist/internal/user"
)

type Container struct {
	AuthHandler  *auth.AuthHandler
	UserHandler  *user.UserHandler
	NurseHandler *nurse.NurseHandler
	AdminHandler *admin.AdminHandler
	ChatHub      *chat.Hub
	ChatHandler  *chat.ChatHandler
}

func NewContainer() *Container {
	// Inicializa o banco de dados
	db := config.GetMongoDB()
	// Construtores: repository → service → handler
	userRepository := repository.NewUserRepository(db)
	nurseRepository := repository.NewNurseRepository(db)
	visitRepository := repository.NewVisitRepository(db)
	messageRepository := repository.NewMessageRepository(db)
	reviewRepository := repository.NewReviewRepository(db)
	hub := chat.NewHub(messageRepository)

	authService := auth.NewAuthService(userRepository, nurseRepository)
	adminService := admin.NewAdminService(userRepository, nurseRepository, visitRepository)
	userService := user.NewUserService(userRepository, nurseRepository, visitRepository, reviewRepository)
	nurseService := nurse.NewNurseService(userRepository, nurseRepository, visitRepository)

	authHandler := auth.NewAuthHandler(authService)
	adminHandler := admin.NewAdminHandler(adminService)
	userHandler := user.NewUserHandler(userService)
	nurseHandler := nurse.NewNurseHandler(nurseService)
	chatHandler := chat.NewChatHandler(messageRepository)

	return &Container{
		AuthHandler:  authHandler,
		AdminHandler: adminHandler,
		UserHandler:  userHandler,
		NurseHandler: nurseHandler,
		ChatHub:      hub,
		ChatHandler:  chatHandler,
	}
}
