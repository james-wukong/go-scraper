package v1

import (
	"net/http"
	"strings"

	"github.com/james-wukong/go-app/cmd/cron/scrapers"
	"github.com/james-wukong/go-app/internal/constants"
	"github.com/james-wukong/go-app/internal/datasources/records"
	V1Postgres "github.com/james-wukong/go-app/internal/datasources/repositories/postgres/v1"
	"github.com/james-wukong/go-app/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/gocolly/colly/v2"
)

func CostcoCategoryScraper(conn *sqlx.DB) {
	// the URL to scrape
	pageDomain := "https://www.costco.com/"
	pageToScrape := pageDomain + "SiteMapDisplayView"

	// creating a new Colly instance
	c := colly.NewCollector()

	// setting a valid User-Agent header
	c.UserAgent = scrapers.CostcoUserAgent
	cookies := []*http.Cookie{
		{
			Name:  scrapers.CostcoCookiesName,
			Value: scrapers.CostcoCookiesValue,
		},
		{
			Name:  "Domain",
			Value: "www.costco.com",
		},
	}
	c.SetCookies(pageToScrape, cookies)

	db := V1Postgres.NewCategoryRepo(conn)

	// initializing the list of pages to scrape with an empty slice
	c.OnHTML(`div.costcoBD-sitemap div.sitemap-section`, func(e *colly.HTMLElement) {
		category := records.Categories{Platform: uint(constants.COSTCO)}
		e.ForEach(`div`, func(_ int, div *colly.HTMLElement) {
			// first level of category
			var level uint = 0
			name := strings.TrimSpace(div.ChildText(`a.h2-style-guide`))
			url := strings.TrimSpace(div.ChildAttr(`a.h2-style-guide`, "href"))
			if len(name) > 0 {
				category.Level = level
				category.Name = name
				category.ParentId = nil
				category.Url = url
				// save category and get categoryId
				categoryDomain := category.ToV1Domain()
				lvl0CategoryId, _ := db.Upsert(&categoryDomain)

				// second level of category
				div.ForEach(`ul:first-of-type > li > ul > li`, func(_ int, lvl1 *colly.HTMLElement) {
					level = 1
					name = strings.TrimSpace(lvl1.ChildText(`a.body-copy-link`))
					url = strings.TrimSpace(lvl1.ChildAttr(`a.body-copy-link`, "href"))
					if len(name) > 0 {
						category.Level = level
						category.Name = name
						category.ParentId = &lvl0CategoryId
						category.Url = url
						// save category and get categoryId
						categoryDomain = category.ToV1Domain()
						lvl1CategoryId, _ := db.Upsert(&categoryDomain)
						// third level of category
						lvl1.ForEach(`ul.sub-list > li`, func(_ int, lvl2 *colly.HTMLElement) {
							level = 2
							name = strings.TrimSpace(lvl2.ChildText(`a`))
							url = strings.TrimSpace(lvl2.ChildAttr(`a`, "href"))
							if len(name) > 0 {
								category.Level = level
								category.Name = name
								category.ParentId = &lvl1CategoryId
								category.Url = url

								// save category and get categoryId
								categoryDomain = category.ToV1Domain()
								db.Upsert(&categoryDomain)
							}
						})
					}
				})
			}
		})
	})

	c.OnScraped(func(r *colly.Response) {
		logger.Debug("c scraper exited ", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryScraper})
	})
	c.Visit(pageToScrape)
}
