package server

import (
	"github.com/m1thrandir225/galore-services/internal/handlers/users"
)

func (s *Server) SetupRoutes(
	userHandler *users.Handler,
) {
	v1 := s.router.Group("/api/v1")

	authRoutes := v1.Group("/")
	authRoutes.Use(authMiddleware(s.tokenMaker))
	userHandler.RegisterRoutes(authRoutes.Group("/users"))
}
