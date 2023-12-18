package pkg

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/config"
	"ElectricCarsServer/ElectricCarsServer/internal/app/handlers"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Application struct {
	Config  *config.Config
	Logger  *logrus.Logger
	Router  *gin.Engine
	Handler *handlers.Handler
}

func NewApp(c *config.Config, r *gin.Engine, l *logrus.Logger, h *handlers.Handler) *Application {
	return &Application{
		Config:  c,
		Logger:  l,
		Router:  r,
		Handler: h,
	}
}

func (a *Application) RunApp() {
	a.Logger.Info("Server start up")
	a.Handler.RegisterHandler(a.Router)

	serverAddress := fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)
	serverAddress = "0.0.0.0:8080" // TODO: fix
	if err := a.Router.Run(serverAddress); err != nil {
		a.Logger.Fatalln(err)
	}
	a.Logger.Info("Server down")
}
