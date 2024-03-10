package api

import (
	db "gopsql/banking/db/sqlc"
	"gopsql/banking/token"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store      db.Store
	tokenMaker token.TokenMaker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and set up routing
func NewServer(store db.Store) *Server {
	tokenMaker, err := token.NewJWTMaker()
	if err != nil {
		log.Fatal("error creating token maker: ", err)
	}

	server := &Server{store: store, tokenMaker: tokenMaker}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	authRoutes := router.Group("/").Use(authMiddleware(tokenMaker))

	// add routes to router
	router.POST("/users", server.CreateUser)
	router.POST("/tokens", server.CreateToken)
	router.POST("/accounts", server.CreateAccount)
	authRoutes.GET("/accounts", server.ListAccounts)
	router.GET("/accounts/:id", server.GetAccount)
	router.POST("/transfers", server.CreateTransfer)
	server.router = router
	return server
}

// Start runs the HTTP server to a specific address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
