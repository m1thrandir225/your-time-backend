package api

import (
	db "m1thrandir225/your_time/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server {
		store: store,
	}

	router := gin.Default();

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// }

	router.POST("/users", server.createUser)
	router.GET("/users/:id", server.getUser)
	server.router = router

	return server;
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}