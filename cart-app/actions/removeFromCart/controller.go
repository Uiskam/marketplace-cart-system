package removeFromCart

import (
	"cart-app/app/common"
	"cart-app/app/cqrs"
	"github.com/gin-gonic/gin"
)

func Controller(ctx *gin.Context) {
	commandBus := ctx.MustGet(common.CommandBusKey).(*cqrs.Bus)

	var cmd Command
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := commandBus.Send(cmd)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, result)
}
