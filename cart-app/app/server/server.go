package server

import (
	"cart-app/actions/addToCart"
	"cart-app/actions/checkoutCart"
	"cart-app/actions/createCart"
	"cart-app/actions/getCart"
	"cart-app/actions/removeFromCart"
	"cart-app/app/common"
	"cart-app/app/cqrs"
	"cart-app/external"
	"cart-app/repository/cart"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"reflect"
)

type Server struct {
	router *gin.Engine
	db     *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	router := gin.Default()

	server := &Server{
		router: router,
		db:     db,
	}
	readRepo := cart.NewReadRepository(db)
	writeRepo := cart.NewWriteRepository(db)

	// Initialize ProductClient for communicating with product-api
	productClient := external.NewProductClient("http://product-app:8080")
	notificationClient := external.NewNotificationClient("http://api:8082")

	queryProcessor := cqrs.NewProcessor()
	queryProcessor.AddHandler(reflect.TypeOf(getCart.Query{}), getCart.NewHandler(readRepo))
	queryBus := cqrs.NewBus(queryProcessor)

	commandProcessor := cqrs.NewProcessor()
	commandProcessor.AddHandler(reflect.TypeOf(createCart.Command{}), createCart.NewHandler(writeRepo))
	commandProcessor.AddHandler(reflect.TypeOf(addToCart.Command{}), addToCart.NewHandler(writeRepo, productClient))
	commandProcessor.AddHandler(reflect.TypeOf(removeFromCart.Command{}), removeFromCart.NewHandler(writeRepo, productClient))
	commandProcessor.AddHandler(reflect.TypeOf(checkoutCart.Command{}), checkoutCart.NewHandler(writeRepo, queryBus, productClient, notificationClient))
	commandBus := cqrs.NewBus(commandProcessor)

	router.Use(func(ctx *gin.Context) {
		ctx.Set(common.QueryBusKey, queryBus)
		ctx.Set(common.CommandBusKey, commandBus)
	})

	server.router.GET("/health", server.healthCheck)
	server.router.POST("/cart/create", createCart.Controller)
	server.router.GET("/cart/:cart_uuid", getCart.Controller)
	server.router.POST("/cart/add", addToCart.Controller)
	server.router.POST("/cart/remove", removeFromCart.Controller)
	server.router.POST("/cart/checkout", checkoutCart.Controller)

	return server
}

func (s *Server) healthCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "ok",
	})
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
