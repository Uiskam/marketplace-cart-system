package server

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"product-app/actions/getProduct"
	"product-app/actions/lockProduct"
	"product-app/actions/sellProduct"
	"product-app/actions/unlockProduct"
	"product-app/app/common"
	"product-app/app/cqrs"
	"product-app/repository/product"
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
	readRepo := product.NewReadRepository(db)
	writeRepo := product.NewWriteRepository(db)

	queryProcessor := cqrs.NewProcessor()
	queryProcessor.AddHandler(reflect.TypeOf(getProduct.Query{}), getProduct.NewHandler(readRepo))
	queryBus := cqrs.NewBus(queryProcessor)

	commandProcessor := cqrs.NewProcessor()
	commandProcessor.AddHandler(reflect.TypeOf(lockProduct.Command{}), lockProduct.NewHandler(writeRepo))
	commandProcessor.AddHandler(reflect.TypeOf(sellProduct.Command{}), sellProduct.NewHandler(writeRepo))
	commandProcessor.AddHandler(reflect.TypeOf(unlockProduct.Command{}), unlockProduct.NewHandler(writeRepo))
	commandBus := cqrs.NewBus(commandProcessor)

	router.Use(func(ctx *gin.Context) {
		ctx.Set(common.QueryBusKey, queryBus)
		ctx.Set(common.CommandBusKey, commandBus)
	})

	server.router.GET("/health", server.healthCheck)
	server.router.GET("/products", getProduct.Controller)
	server.router.POST("/products/lock", lockProduct.Controller)
	server.router.POST("/products/sell", sellProduct.Controller)
	server.router.POST("/products/unlock", unlockProduct.Controller)

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
