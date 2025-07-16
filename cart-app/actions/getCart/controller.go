package getCart

import (
	"cart-app/app/common"
	"cart-app/app/cqrs"
	"github.com/gin-gonic/gin"
)

func Controller(ctx *gin.Context) {
	queryBus := ctx.MustGet(common.QueryBusKey).(*cqrs.Bus)

	cartUUID := ctx.Param("cart_uuid")
	query := Query{CartUUID: cartUUID}

	result, err := queryBus.Send(query)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, result)
}
