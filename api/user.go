package api

import (
	"database/sql"
	db "m1thrandir225/your_time/db/sqlc"
	"m1thrandir225/your_time/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type createUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type createUserResponse struct {
	ID uuid.UUID `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type getUserRequest struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type getUserResponse struct { 
	ID uuid.UUID `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdateAt string `json:"updated_at"`
}


func (server *Server) createUser(ctx * gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams {
		FirstName: req.FirstName,
		LastName: req.LastName,
		Email: req.Email,
		Password: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	responseData := createUserResponse {
		FirstName: user.FirstName,
		ID: user.ID,
		LastName: user.LastName,
		Email: user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	ctx.JSON(http.StatusOK, responseData)

}


func (server *Server) getUser(ctx * gin.Context) {
	var req getUserRequest;

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return	
	}
	user, err := server.store.GetUserByID(ctx, uuid.MustParse(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	responseData := getUserResponse {
		ID: user.ID,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdateAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	ctx.JSON(http.StatusOK, responseData)
}