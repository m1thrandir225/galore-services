package shared

import "github.com/gin-gonic/gin"

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func MessageResponse(message string) gin.H {
	return gin.H{"message": message}
}
