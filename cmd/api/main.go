package main

import (
	"runtime"

	"github.com/james-wukong/go-app/cmd/api/server"
	"github.com/james-wukong/go-app/internal/config"
	"github.com/james-wukong/go-app/internal/constants"
	"github.com/james-wukong/go-app/pkg/logger"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := config.InitializeAppConfig(); err != nil {
		logger.Fatal(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig})
	}
	logger.Info("configuration loaded", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig})
}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9001
// @BasePath  /api/

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	numCPU := runtime.NumCPU()
	logger.InfoF("The project is running on %d CPU(s)", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig}, numCPU)

	if runtime.NumCPU() > 2 {
		runtime.GOMAXPROCS(numCPU / 2)
	}

	app, err := server.NewApp()
	if err != nil {
		logger.Panic(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryServer})
	}
	if err := app.Run(); err != nil {
		logger.Fatal(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryServer})
	}
}
