package handler

import (
	_ "RIP/docs"
	"RIP/internal/app/config"
	"RIP/internal/app/pkg/hash"
	"RIP/internal/app/redis"
	"RIP/internal/app/repository"
	"RIP/internal/app/role"
	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
)

const (
	creatorID   = 1
	moderatorID = 1
)

type Handler struct {
	Logger     *logrus.Logger
	Repository *repository.Repository
	Minio      *minio.Client
	Config     *config.Config
	Redis      *redis.Client
	//TokenManager auth.TokenManager
	Hasher hash.PasswordHasher
}

func NewHandler(
	l *logrus.Logger,
	r *repository.Repository,
	m *minio.Client,
	conf *config.Config,
	red *redis.Client,
	// tokenManager auth.TokenManager,
) *Handler {
	return &Handler{
		Logger:     l,
		Repository: r,
		Minio:      m,
		Config:     conf,
		Redis:      red,
		//TokenManager: tokenManager,
		Hasher: hash.NewSHA256Hasher(os.Getenv("SALT")),
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	// услуги
	api.GET("/autoparts", h.WithoutJWTError(role.Buyer, role.Moderator, role.Admin), h.AutopartsList) // ?
	api.GET("/autoparts/:id", h.GetAutopartById)                                                      // ?
	api.POST("/autoparts", h.WithAuthCheck(role.Moderator, role.Admin), h.AddAutopart)
	api.PUT("/autoparts", h.WithAuthCheck(role.Moderator, role.Admin), h.UpdateAutopart)
	api.PUT("/autoparts/upload-image", h.WithAuthCheck(role.Moderator, role.Admin), h.AddImage)
	api.DELETE("/autoparts", h.WithAuthCheck(role.Moderator, role.Admin), h.DeleteAutopart)
	api.POST("/autoparts/request", h.WithAuthCheck(role.Buyer, role.Moderator, role.Admin), h.AddAutopartToRequest)
	api.Use(cors.Default()).DELETE("/autoparts/delete/:id", h.DeleteAutopart)

	// заявки
	api.GET("/assemblies", h.WithAuthCheck(role.Buyer, role.Moderator, role.Admin), h.AssemblyList)
	api.GET("/assemblies/:id", h.WithAuthCheck(role.Buyer, role.Moderator, role.Admin), h.GetAssemblyById)
	api.GET("/assemblies/current", h.WithAuthCheck(role.Buyer, role.Moderator, role.Admin), h.AssemblyCurrent)
	api.PUT("/assemblies", h.WithAuthCheck(role.Buyer, role.Moderator, role.Admin), h.UpdateAssembly)

	// статусы
	api.PUT("/assemblies/form", h.WithAuthCheck(role.Buyer, role.Moderator, role.Admin), h.FormAssemblyRequest)
	api.PUT("/assemblies/updateStatus", h.WithAuthCheck(role.Moderator, role.Admin), h.UpdateStatusAssemblyRequest)

	api.DELETE("/assemblies", h.WithAuthCheck(role.Buyer, role.Moderator, role.Admin), h.DeleteAssembly)

	// m-m
	api.DELETE("/assembly-request-autopart", h.WithoutJWTError(role.Buyer, role.Moderator, role.Admin), h.DeleteAutopartFromRequest)
	api.PUT("/assembly-request-autopart", h.WithoutJWTError(role.Buyer, role.Moderator, role.Admin), h.UpdateAssemblyAutopart)
	registerStatic(router)

	// auth && reg
	api.POST("/user/signIn", h.Login)
	api.POST("/user/signUp", h.Register)
	api.POST("/user/logout", h.Logout)

	// асинхронный сервис
	api.PUT("/assemblies/user-form-start", h.WithoutJWTError(role.Buyer, role.Moderator, role.Admin), h.UserRequest) // обращение к асинхронному сервису
	api.PUT("/assemblies/user-form-finish", h.FinishUserRequest)                                                     // обращение к асинхронному сервису

}

func registerStatic(router *gin.Engine) {
	router.LoadHTMLGlob("static/templates/*")
	router.Static("/static", "./static")
	router.Static("/css", "./static")
	router.Static("/img", "./static")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	h.Logger.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      errorStatusCode,
		"description": err.Error(),
	})
}

func (h *Handler) successHandler(ctx *gin.Context, key string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		key:      data,
	})
}

func (h *Handler) successAddHandler(ctx *gin.Context, key string, data interface{}) {
	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		key:      data,
	})
}
