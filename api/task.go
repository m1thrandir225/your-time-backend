package api

import (
	"database/sql"
	db "m1thrandir225/your_time/db/sqlc"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type createTaskRequest struct {
	Title string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	DueDate string `json:"due_date" binding:"required"`
	ReminderDate *string `json:"reminder_date"`
	UserID string `json:"user_id" binding:"required"`
}

type createTaskResponse struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	DueDate string `json:"due_date"`
	ReminderDate string `json:"reminder_date"`
	UserID string `json:"user_id"`
}

type getTaskByIDRequest struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type getTaskByUser struct {
	UserID string `uri:"user_id" binding:"required,min=1"`
}


func (server *Server) createTask(ctx *gin.Context) {
	var req createTaskRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err));
		return;
	}

	dueDate, err := time.Parse(time.RFC3339, req.DueDate);

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err));
		return;
	}

	reminderDate, err := time.Parse(time.RFC3339, *req.ReminderDate);

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err));
		return;
	}

	userUUID, err := uuid.Parse(req.UserID);

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err));
		return;
	}


	arg := db.CreateTaskParams {
		Title: req.Title,
		Description: sql.NullString{String: req.Description, Valid: true},
		DueDate: dueDate,
		ReminderDate: sql.NullTime{Time: reminderDate, Valid: true},
		UserID: userUUID,
	}

	task, err := server.store.CreateTask(ctx, arg);

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err));
		return;
	}

	res := createTaskResponse {
		ID: task.ID.String(),
		Title: task.Title,
		Description: task.Description.String,
		DueDate: task.DueDate.Format(time.RFC3339),
		ReminderDate: task.ReminderDate.Time.Format(time.RFC3339),
		UserID: task.UserID.String(),
	}
	ctx.JSON(http.StatusOK, res);
}