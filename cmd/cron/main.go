package main

import (
	"flag"

	V1Costco "github.com/james-wukong/go-app/cmd/cron/scrapers/costco/v1"
	V1FBasics "github.com/james-wukong/go-app/cmd/cron/scrapers/foodbasics/v1"
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

const (
	defaultCmd    = "discounts"
	discCostcoCmd = "discount-costco"
	categoryCmd   = "category"
	allCmd        = "all"
)

var (
	strCmd *string
)

func main() {
	// setup databases
	conn, err := utils.SetupPostgresConnection()
	if err != nil {
		logger.Info("error caught: "+err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryScraper})
	}
	defer conn.Close()

	strCmd = flag.String("run", defaultCmd, "execute cmd: 1. category 2. discounts 3. all")
	flag.Parse()
	switch *strCmd {
	case defaultCmd:
		V1Costco.CostcoScraper(conn)
		V1FBasics.FBasicsScraper(conn)
	case discCostcoCmd:
		V1Costco.CostcoScraper(conn)
	case categoryCmd:
		V1Costco.CostcoCategoryScraper(conn)
		V1FBasics.FBasicsCategoryScraper(conn)
	case allCmd:
		V1Costco.CostcoScraper(conn)
		V1FBasics.FBasicsScraper(conn)
		V1Costco.CostcoCategoryScraper(conn)
		V1FBasics.FBasicsCategoryScraper(conn)
	default:
		logger.Info("usage: go run cmd/cron/main.go -run=category|discounts|all", logrus.Fields{"-run": *strCmd})
	}

}
