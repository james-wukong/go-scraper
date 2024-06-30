package main

import (
	V1Costco "github.com/james-wukong/go-app/cmd/cron/scrapers/costco/v1"
	"github.com/james-wukong/go-app/internal/config"
	"github.com/james-wukong/go-app/internal/constants"
	"github.com/james-wukong/go-app/internal/utils"
	"github.com/james-wukong/go-app/pkg/logger"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := config.InitializeAppConfig(); err != nil {
		logger.Fatal(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig})
	}
	logger.Info("configuration loaded", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig})
}

func main() {
	// setup databases
	conn, err := utils.SetupPostgresConnection()
	if err != nil {
		logger.Info("error caught: "+err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryScraper})
	}
	defer conn.Close()
	// scraper categories of costco
	// V1Costco.CostcoCategoryScrap(conn)
	// scraper promotion products
	V1Costco.CostcoScrap(conn)
}
