package router

import (
	"medassist/internal/di"
	"medassist/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(r *gin.RouterGroup, container *di.Container) {	
	admin := r.Group("/admin")
	{
		//CRUS DE USER/ NURSE  TODO
		admin.GET("/dashboard", middleware.AuthAdmin(),container.AdminHandler.AdminDashboard)
		admin.GET("/all_pending_registers", container.AdminHandler.GetRegistersToApprove)
		admin.GET("/documents/:id", middleware.AuthAdmin(), container.AdminHandler.GetDocuments)
		admin.GET("/download/:id", container.AdminHandler.DownloadFile)
		admin.PATCH("/approve/:id", middleware.AuthAdmin(), container.AdminHandler.ApproveNurseRegister)
		admin.POST("/reject/:id", middleware.AuthAdmin(), container.AdminHandler.RejectNurseRegister)
		admin.GET("/file/:id", container.UserHandler.GetFileByID)
		admin.GET("/users", middleware.AuthAdmin(), container.AdminHandler.UsersManagement)
		admin.PATCH("/user/:id", middleware.AuthAdmin(), container.AdminHandler.UpdateUser)
		admin.DELETE("/user/:id", middleware.AuthAdmin(), container.AdminHandler.DeleteUser)

		//crud visits

	}
}