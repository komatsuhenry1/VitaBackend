package di

import (
	"medassist/config"
	"medassist/internal/admin"
	"medassist/internal/auth"
	"medassist/internal/chat"
	"medassist/internal/nurse"
	"medassist/internal/repository"
	"medassist/internal/user"
	"medassist/internal/payment"
)

type Container struct {
	AuthHandler  *auth.AuthHandler
	UserHandler  *user.UserHandler
	NurseHandler *nurse.NurseHandler
	AdminHandler *admin.AdminHandler
	ChatHub      *chat.Hub
	ChatHandler  *chat.ChatHandler
	PaymentHandler *payment.PaymentHandler
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
	paymentRepository := repository.NewPaymentRepository(db)
	stripeRepository := repository.NewStripeRepository()
	hub := chat.NewHub(messageRepository)

	authService := auth.NewAuthService(userRepository, nurseRepository)
	adminService := admin.NewAdminService(userRepository, nurseRepository, visitRepository)
	userService := user.NewUserService(userRepository, nurseRepository, visitRepository, reviewRepository, hub)
	nurseService := nurse.NewNurseService(userRepository, nurseRepository, visitRepository, reviewRepository, stripeRepository)
	paymentService := payment.NewPaymentService(paymentRepository, userRepository)


	authHandler := auth.NewAuthHandler(authService)
	adminHandler := admin.NewAdminHandler(adminService)
	userHandler := user.NewUserHandler(userService)
	nurseHandler := nurse.NewNurseHandler(nurseService)
	chatHandler := chat.NewChatHandler(messageRepository)
	paymentHandler := payment.NewPaymentHandler(paymentService)

	return &Container{
		AuthHandler:  authHandler,
		AdminHandler: adminHandler,
		UserHandler:  userHandler,
		NurseHandler: nurseHandler,
		ChatHub:      hub,
		ChatHandler:  chatHandler,
		PaymentHandler: paymentHandler,
	}
}
