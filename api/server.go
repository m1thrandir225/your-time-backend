package api

import (
	"fmt"
	db "m1thrandir225/your_time/db/sqlc"
	token "m1thrandir225/your_time/token"
	"m1thrandir225/your_time/util"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config util.Config
	store db.Store
	tokenMaker token.Maker
	router *gin.Engine
}

func NewServer(config util.Config, store db.Store)( *Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server {
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}

	server.SetupRouter()

	
	return server, nil;
}

func (server *Server) SetupRouter() {
	router := gin.Default();

	//Register
	router.POST("/users", server.createUser)

	//Login
	router.POST("/users/login", server.loginUser)

	router.POST("/tasks", server.createTask)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/users/:id", server.getUser)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}