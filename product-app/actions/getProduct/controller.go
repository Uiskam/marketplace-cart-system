package getProduct

import (
	"github.com/gin-gonic/gin"
	"product-app/app/common"
	"product-app/app/cqrs"
)

func Controller(ctx *gin.Context) {
	queryBus := ctx.MustGet(common.QueryBusKey).(*cqrs.Bus)
	result, err := queryBus.Send(Query{})
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
	}
	ctx.JSON(200, result)
	return
}
