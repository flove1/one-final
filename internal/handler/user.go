package handler

import (
	"errors"
	"net/http"
	"one-lab-final/internal/entity"
	"one-lab-final/internal/handler/api"
	"one-lab-final/pkg/util"

	"github.com/gin-gonic/gin"
)

// var passwordRegEx = regexp.MustCompile(`^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{4,8}$`)

// @Summary      Register new user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param data body api.CreateUserRequest true "Request body"
//
// @Success      201 {object} api.DefaultResponse "User succesfully created"
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /user/register [post]
func (h *Handler) createUser(ctx *gin.Context) {
	var req api.CreateUserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &api.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// if !passwordRegEx.Match([]byte(req.Password)) {
	// 	ctx.JSON(http.StatusBadRequest, &ErrorResponse{
	// 		http.StatusBadRequest,
	// 		"password is not strong enough",
	// 	})
	// }

	err = h.Services.CreateUser(ctx, &entity.User{
		Username:  &req.Username,
		Email:     &req.Email,
		FirstName: &req.FirstName,
		LastName:  &req.LastName,
		Password: entity.Password{
			Plaintext: &req.Password,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &api.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, &api.DefaultResponse{
		Code:    http.StatusCreated,
		Message: "user succesfully created",
	})
}

// @Summary      Authenticated user and return access token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param data body api.LoginRequest true "Request body"
//
// @Success      201 {object} api.LoginResponse "token succesfully created"
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /user/login [post]
func (h *Handler) login(ctx *gin.Context) {
	var req api.LoginRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &api.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	token, err := h.Services.Login(ctx, req.Credentials, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, util.ErrMismatchedPassword):
			ctx.JSON(http.StatusBadRequest, &api.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, &api.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, &api.LoginResponse{
		Code:    http.StatusCreated,
		Message: "token succesfully created",
		Body:    token,
	})
}

// @Summary      Update user info
// @Tags         Users
// @Produce      json
// @Security ApiKeyAuth
// @Param data body api.LoginRequest true "Request body"
//
// @Success      201 {object} api.LoginResponse "token succesfully created"
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /user/update [patch]
func (h *Handler) updateUser(ctx *gin.Context) {
	var req api.UpdateUserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &api.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// if !passwordRegEx.Match([]byte(req.Password)) {
	// 	ctx.JSON(http.StatusBadRequest, &ErrorResponse{
	// 		http.StatusBadRequest,
	// 		"password is not strong enough",
	// 	})
	// }

	userID := ctx.MustGet("userID").(int64)

	err = h.Services.UpdateUser(ctx, &entity.User{
		ID:        userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password: entity.Password{
			Plaintext: req.Password,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &api.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &api.DefaultResponse{
		Code:    http.StatusOK,
		Message: "user succesfully updated",
	})
}

// @Summary      Delete user
// @Tags         Users
// @Produce      json
// @Security ApiKeyAuth
//
// @Success      200 {object} api.DefaultResponse "user succesfully deleted"
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /user/delete [delete]
func (h *Handler) deleteUser(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(int64)

	err := h.Services.DeleteUser(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &api.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &api.DefaultResponse{
		Code:    http.StatusOK,
		Message: "user succesfully deleted",
	})
}
