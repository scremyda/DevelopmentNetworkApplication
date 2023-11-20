package handlers

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/ds"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) AddUser(ctx *gin.Context) {
	var newUser ds.Users
	if err := ctx.BindJSON(&newUser); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if newUser.Login == "" {
		h.errorHandler(ctx, http.StatusBadRequest, loginCantBeEmpty)
		return
	}
	if newUser.Password == "" {
		h.errorHandler(ctx, http.StatusBadRequest, passwordCantBeEmpty)
		return
	}
	if err := h.Repository.AddUser(&newUser); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.successAddHandler(ctx, "user_created", gin.H{
		"user_id":      newUser.ID,
		"login":        newUser.Login,
		"is_moderator": newUser.IsModerator,
	})

	// h.successHandler(ctx, "user_created", newUser)
}
