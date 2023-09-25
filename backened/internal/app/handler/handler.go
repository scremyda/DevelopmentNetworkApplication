package handler

import (
	"backened/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Logger     *logrus.Logger
	Repository *repository.Repository
}

func NewHandler(l *logrus.Logger, r *repository.Repository) *Handler {
	return &Handler{
		Logger:     l,
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/", h.AutopartsList)
	router.GET("/autoparts/:id", h.AutopartById)
	router.POST("/delete/:id", h.Deleteautopart)
	registerStatic(router)
}

func registerStatic(router *gin.Engine) {
	router.LoadHTMLGlob("backened/static/templates/*")
	router.Static("/static", "././backened/static")
	router.Static("/css", "./static")
	router.Static("/img", "./static")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	h.Logger.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
