package handlers

import (
	_ "ElectricCarsServer/ElectricCarsServer/docs"
	"ElectricCarsServer/ElectricCarsServer/internal/app/config"
	models "ElectricCarsServer/ElectricCarsServer/internal/app/ds"
	"ElectricCarsServer/ElectricCarsServer/internal/app/pkg/auth"
	"ElectricCarsServer/ElectricCarsServer/internal/app/pkg/hash"
	"ElectricCarsServer/ElectricCarsServer/internal/app/redis"
	"ElectricCarsServer/ElectricCarsServer/internal/app/repo"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
)

const (
	baseURL = "api"

	autoparts        = baseURL + "/autoparts"
	autopartsID      = baseURL + "/autoparts/:id"
	autopartsList    = baseURL + "/autoparts/get-all"
	addAutopartImage = baseURL + "/autoparts/upload-image"
	addAssembly      = baseURL + "/autoparts/add-to-assembly"

	assembly         = baseURL + "/assembly"
	assemblyForm     = baseURL + "/assembly/form"
	assemblyComplete = baseURL + "/assembly/complete"
	assemblyReject   = baseURL + "/assembly/reject"
	assemblyID       = baseURL + "/assembly/:id"
	assemblyList     = baseURL + "/assembly/get-all"

	autoparts_assembly = baseURL + "/autoparts_assembly"

	user = baseURL + "/user"
)

type Handler struct {
	Logger       *logrus.Logger
	Repository   *repo.Repository
	Minio        *minio.Client
	Config       *config.Config
	Redis        *redis.Client
	TokenManager auth.TokenManager
	Hasher       hash.PasswordHasher
}

func NewHandler(
	l *logrus.Logger,

	r *repo.Repository,
	m *minio.Client,
	conf *config.Config,
	red *redis.Client,

	tokenManager auth.TokenManager,
) *Handler {
	return &Handler{
		Logger:       l,
		Repository:   r,
		Minio:        m,
		Config:       conf,
		Redis:        red,
		TokenManager: tokenManager,
		Hasher:       hash.NewSHA256Hasher(os.Getenv("SALT")),
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {

	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET(autopartsList, h.AutopartsList)
	router.GET(autopartsID, h.AutopartById)

	router.POST(autoparts, h.WithAuthCheck([]models.Role{models.Admin}), h.AddAutopart)
	router.POST(addAutopartImage, h.WithAuthCheck([]models.Role{models.Admin}), h.AddImage)
	router.POST(addAssembly, h.WithAuthCheck([]models.Role{models.Client}), h.AddToAssembly)

	router.PUT(autoparts, h.WithAuthCheck([]models.Role{models.Admin}), h.UpdateAutopart)

	router.DELETE(autopartsID, h.WithAuthCheck([]models.Role{models.Admin}), h.DeleteAutopart)
	//=============================================//

	router.GET(assemblyList, h.WithAuthCheck([]models.Role{models.Admin, models.Client}), h.AssembliesList)

	// разный доступ, у админа к любой, у юзера только к своей
	router.GET(assemblyID, h.WithAuthCheck([]models.Role{models.Client, models.Admin}), h.AssemblyById)

	router.PUT(assembly, h.WithAuthCheck([]models.Role{models.Admin}), h.UpdateAssembly)

	// статусы
	router.PUT(assemblyForm, h.WithAuthCheck([]models.Role{models.Client}), h.FormAssembly)
	router.PUT(assemblyComplete, h.WithAuthCheck([]models.Role{models.Admin}), h.CompleteAssembly)
	router.PUT(assemblyReject, h.WithAuthCheck([]models.Role{models.Admin}), h.RejectAssembly)

	router.DELETE(assembly, h.WithAuthCheck([]models.Role{models.Client}), h.DeleteAssembly)
	//=============================================//

	router.PUT(autoparts_assembly, h.UpdateCountAutopartAssembly)

	router.DELETE(autoparts_assembly, h.DeleteFromAssembly)
	//=============================================//

	router.PUT(user, h.AddUser)
	// авторизация и регистрация
	router.POST(user+"/signIn", h.SignIn)
	router.POST(user+"/signUp", h.SignUp)
	router.POST(user+"/logout", h.Logout)

	//=============================================//
	registerStatic(router)
}

func registerStatic(router *gin.Engine) {
	//router.LoadHTMLGlob("ElectricCarsServer/static/templates/*")
	//router.Static("/static", "././backened/static")
	//router.Static("/css", "./static")
	//router.Static("/img", "./static")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	h.Logger.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
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
