package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendNotification(ctx *gin.Context) {
	service := ctx.MustGet(ServiceContextKey).(*Service)
	var notification Notification
	if err := ctx.ShouldBindJSON(&notification); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.EnqueueNotification(&notification); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue notification"})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"message": "Notification enqueued successfully"})
}

func GetPendingTasks(ctx *gin.Context) {
	service := ctx.MustGet(ServiceContextKey).(*Service)
	tasks, err := service.GetPendingTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pending tasks"})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

func GetCompletedTasks(ctx *gin.Context) {
	service := ctx.MustGet(ServiceContextKey).(*Service)
	tasks, err := service.GetCompletedTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sent tasks"})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

func GetFailedTasks(ctx *gin.Context) {
	service := ctx.MustGet(ServiceContextKey).(*Service)
	tasks, err := service.GetFailedTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get failed tasks"})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

func SendNow(ctx *gin.Context) {
	service := ctx.MustGet(ServiceContextKey).(*Service)
	taskID := ctx.Param("id")

	isFound, err := service.SendNow(taskID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send task now"})
		return
	}
	if !isFound {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task sent now successfully"})
}

func CancelTask(ctx *gin.Context) {
	service := ctx.MustGet(ServiceContextKey).(*Service)
	taskID := ctx.Param("id")

	isFound, err := service.CancelTask(taskID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel task"})
		return
	}
	if !isFound {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task cancelled successfully"})
}
