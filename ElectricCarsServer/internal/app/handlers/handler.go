package handlers

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/repo"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	baseURL = "api"

	autoparts        = baseURL + "/autoparts"
	autopartsID      = baseURL + "/autoparts/:id"
	autopartsList    = baseURL + "/autoparts/get-all"
	addAutopartImage = baseURL + "/autoparts/upload-image"
	addAssembly      = baseURL + "/autoparts/add-to-assembly"

	assembly               = baseURL + "/assembly"
	assemblyForm           = baseURL + "/assembly/form"
	assemblyCompleteReject = baseURL + "/assembly/complete_reject"
	assemblyReject         = baseURL + "/assembly/reject"
	assemblyID             = baseURL + "/assembly/:id"
	assemblyList           = baseURL + "/assembly/get-all"

	autoparts_assembly = baseURL + "/autoparts_assembly"

	user = baseURL + "/user"
)

type Handler struct {
	Logger     *logrus.Logger
	Repository *repo.Repository
	Minio      *minio.Client
}

func NewHandler(l *logrus.Logger, r *repo.Repository, m *minio.Client) *Handler {
	return &Handler{
		Logger:     l,
		Repository: r,
		Minio:      m,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {

	router.GET(autopartsList, h.AutopartsList)
	router.GET(autopartsID, h.AutopartById)

	router.POST(autoparts, h.AddAutopart)
	router.POST(addAutopartImage, h.AddImage)
	router.POST(addAssembly, h.AddToAssembly)

	router.PUT(autoparts, h.UpdateAutopart)
	router.PUT(addAutopartImage, h.AddImage)

	router.DELETE(autopartsID, h.DeleteAutopart)
	//=============================================//

	router.GET(assemblyList, h.AssembliesList)
	router.GET(assemblyID, h.AssemblyById)

	router.PUT(assembly, h.UpdateAssembly)
	router.PUT(assemblyForm, h.FormAssembly)
	router.PUT(assemblyCompleteReject, h.CompleteRejectAssembly)
	//router.PUT(assemblyReject, h.RejectAssembly)

	router.DELETE(assembly, h.DeleteAssembly)
	//=============================================//

	router.PUT(autoparts_assembly, h.UpdateCountAutopartAssembly)

	router.DELETE(autoparts_assembly, h.DeleteFromAssembly)
	//=============================================//

	router.PUT(user, h.AddUser)

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
