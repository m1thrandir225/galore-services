package users

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/:id", h.GetUserDetails)
	rg.PUT("/:id", h.UpdateUserInformation)
	rg.DELETE("/:id", h.DeleteUser)
	rg.PUT("/:id/password", h.UpdateUserPassword)
	rg.PUT("/:id/push-notifications", h.UpdateUserPushNotificationsSettings)
	rg.PUT("/:id/email-notifications", h.UpdateUserEmailNotificationsSettings)
}
