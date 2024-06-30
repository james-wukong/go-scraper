package v1

import (
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/james-wukong/go-app/cmd/cron/scrapers"
	"github.com/james-wukong/go-app/internal/constants"
	"github.com/james-wukong/go-app/internal/datasources/records"
	V1Postgres "github.com/james-wukong/go-app/internal/datasources/repositories/postgres/v1"
	"github.com/james-wukong/go-app/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func CostcoScrap(conn *sqlx.DB) {
	// initializing the list of pages to scrape with an empty slice
	var pagesToScrape []string
	var searchPattern string = "/CatalogSearch?"

	// the first pagination URL to scrape
	var pageDomain string = "https://www.costco.ca/"
	pageToScrape := pageDomain + "coupons.html"

	// creating a new Colly instance
	c := colly.NewCollector(
		// colly.Debugger(&debug.LogDebugger{}),
		// colly.AllowURLRevisit(),
		// turning on the asynchronous request mode in Colly
		colly.Async(true),
	)

	// setting a valid User-Agent header
	c.UserAgent = scrapers.CostcoUserAgent
	cookies := []*http.Cookie{
		{
			Name:  scrapers.CostcoCookiesName,
			Value: scrapers.CostcoCookiesValue,
		},
		{
			Name:  "Domain",
			Value: "www.costco.ca",
		},
		{
			Name:  "Domain",
			Value: ".costco.ca",
		},
	}
	c.SetCookies(pageToScrape, cookies)
	// scraper for details
	d := c.Clone()
	var categoryId int
	var prodId string
	// Channel to signal completion of the first callback
	var done chan struct{}
	var valid scrapers.Valid
	var price = make(map[string]scrapers.Prices)

	db := V1Postgres.NewCategoryRepo(conn)
	dbProduct := V1Postgres.NewProductRepo(conn)
	dbAvgRate := V1Postgres.NewAvgRatingRepo(conn)
	dbDiscHist := V1Postgres.NewDiscHistoryRepo(conn)

	// catch the promotion container from coupon link
	c.OnHTML("ul.CLPcontainer", func(e *colly.HTMLElement) {
		// discovering a new page
		e.ForEach(`li.couponbox`, func(i int, p *colly.HTMLElement) {
			link := p.ChildAttr(`a.btn-eOffer`, "href")
			tmp := p.ChildText(`div.productDetails span.sku`)
			// tmp = strings.TrimSpace(strings.Replace(tmp, "Item number", "", -1))
			re := regexp.MustCompile(`[^0-9,]+`)
			tmp = re.ReplaceAllString(tmp, "")

			sku := getIDFromUrl(link)
			if strings.Contains(link, pageDomain) {
				prices := scrapers.Prices{EcoFee: 0.0}
				var txt, pri string
				var validSelector string = `div.CLpbulkcoup > span.CLP-validdates`
				from := p.ChildAttr(validSelector+`> time:first-child`, "datetime")
				to := p.ChildAttr(validSelector+`> time:nth-child(2)`, "datetime")
				duration := p.ChildText(validSelector + `> span`)
				valid.StartAt = parseDate(from)
				valid.EndedAt = parseDate(to)
				valid.Duration = strings.TrimSpace(duration)
				prices.Valid = &valid
				p.ForEach(`div.CLpbulkcoup > div.CLP-product table > tbody > tr`, func(_ int, q *colly.HTMLElement) {
					txt = strings.TrimSpace(q.ChildText(`span.eco-priceTableText`))
					pri = q.ChildText(`span.eco-priceTable`)
					if _, ok := price[sku]; !ok {
						val, err := parsePrice(pri)
						if err == nil {
							switch txt {
							case "In-warehouse":
								prices.InWarehouse = float32(val)
							case "Eco fee":
								prices.EcoFee = float32(val)
							case "Instant savings":
								prices.InstSave = float32(val)
							case "PRICE":
								prices.Price = float32(val)
								// default:
								// 	logger.Debug("txt caught is: ", , logrus.Fields{})
							}
						}
					}
				})
				if strings.Contains(link, searchPattern) {
					// if this is a search result link
					if tmpSlice := strings.Split(tmp, ","); len(tmpSlice) > 0 {
						for _, itemSku := range tmpSlice {
							price[itemSku] = prices
						}
					}
					p.Request.Visit(link)
				} else if sku != "" {
					// if this is a product detail link
					price[sku] = prices
					pagesToScrape = append(pagesToScrape, link)
				}
			} else {
				logger.Debug("empty link1: "+link, logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryScraper})
			}
		})
	})

	// catch the catelog search result container
	c.OnHTML(`div.product-list`, func(e *colly.HTMLElement) {
		// discovering a new page
		e.ForEach(`div.product`, func(_ int, p *colly.HTMLElement) {
			link := p.ChildAttr(`span.description a`, "href")
			// itemSku := p.ChildAttr(`div.product-tile-set > input.itemNumber[type="hidden"]`, "value")
			// // sku := getIDFromUrl(link)
			// pAfterStr := p.ChildText(`div.product-tile-set div.thumbnail div.price`)
			// pAfter, err := parsePrice(pAfterStr)
			// if err != nil {
			// 	pAfter = 0.0
			// }
			// pSaveStr := p.ChildText(`div.product-tile-set div.thumbnail p.promo`)
			// pSave, err := parsePrice(pSaveStr)
			// if err != nil {
			// 	pSave = 0.0
			// }
			// prices := scrapers.Prices{
			// 	InWarehouse: float32(pAfter + pSave),
			// 	EcoFee:      0.0,
			// 	InstSave:    float32(pSave),
			// 	Price:       float32(pAfter),
			// 	Valid:       nil,
			// }
			// price[itemSku] = prices
			if strings.Contains(p.Request.URL.String(), searchPattern) {
				if strings.Contains(link, pageDomain) {
					// if this is a product detail link
					pagesToScrape = append(pagesToScrape, link)
				}
			} else {
				logger.Debug("empty link2: "+link, logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryScraper})
			}
		})
	})

	c.OnScraped(func(r *colly.Response) {
		if len(pagesToScrape) > 0 {
			logger.Debug("c scraper done, with total length: "+strconv.Itoa(len(pagesToScrape)), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryScraper})
		} else {
			logger.Debug("c scraper not exit properly ", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryScraper})
		}

	})

	// scrapping the catgories
	d.OnHTML(`ul#crumbs_ul > div > ul > li:last-child a[itemprop="item"]`, func(e *colly.HTMLElement) {
		// logger.Debug("sending signal chan: ", logrus.Fields{})
		// discovering a new page
		// category := e.DOM.Contents().Not("i, span").Text()
		// category = strings.TrimSpace(category)
		pathName := getPathFromURL(e.Attr(`href`))
		// categoryD, err := db.GetByNamePlatform(category, uint(constants.COSTCO))
		categoryD, err := db.GetByURLPlatform(pathName, uint(constants.COSTCO))
		categoryId = categoryD.Id
		if err != nil {
			logger.Debug("category error caught: "+err.Error(), logrus.Fields{"category": pathName, "int": uint(constants.COSTCO)})
			categoryId = scrapers.CostcoDefaultCategory
		}
		prodId = getIDFromUrl(e.Request.URL.String())
		close(done) // Signal that the first callback is done
	})

	// scraping promotion details
	d.OnHTML(`div[itemtype="https://schema.org/Product"]`, func(e *colly.HTMLElement) {
		// logger.Debug("receiving signal chan: ", logrus.Fields{})
		<-done // Wait for the signal from the first callback
		var reviews reviewResp
		product := records.Products{}
		details := records.DetailBase{Detail: make(map[string][]string)}
		specs := records.SpecBase{Spec: make(map[string]string)}
		time := records.TimeBase{
			CreatedAt: time.Now(),
		}

		pName := e.ChildText(`div#product-details div.product-h1-container-v2 h1[itemprop="name"]`)
		product.Name = strings.Split(pName, ", ")[0]
		product.CategoryId = categoryId
		product.ProdId = prodId
		itemSku := getIDFromUrl(e.Request.URL.String())
		// product.Sku = getIDFromUrl(e.Request.URL.String())
		product.Sku = e.ChildText(`div#product-details div#product-body-item-number span`)

		var discHist records.DiscountHistories
		var p scrapers.Prices
		var ok bool
		// Check if product.Sku exists in the price map
		if p, ok = price[itemSku]; !ok {
			// If product.Sku does not exist, check if itemSku exists in the price map
			if p, ok = price[product.Sku]; !ok {
				// If neither exists, skip the rest of the code
				logger.Debug("price & sku not found: ", logrus.Fields{"sku": product.Sku, "price": p})
			}
		}
		product.Model = e.ChildText(`div#product-details div#product-body-model-number span`)
		product.Price = p.InWarehouse
		product.Source = int(constants.COSTCO)
		product.UrlLink = e.Request.URL.String()
		product.ImageLink = e.ChildAttr(`div#left-side-content > div#zoomViewer > div#productImageContainer > div#productImageOverlay > img#initialProductImage`, "src")

		descNode := e.DOM.Find(`div#product-info > div > div > div > div > div#nav-pdp-tab-header-3 > div > div.product-info-description > span#productDescriptions1`)
		description := strings.TrimSpace(descNode.Contents().First().Text())
		details.Detail["Description"] = []string{description}
		// TODO add more details into struct
		// next := descNode.Text()
		e.ForEach(`div#product-info > div > div > div > div > div#nav-pdp-tab-header-5 > div > div.product-info-description > div > div`, func(_ int, s *colly.HTMLElement) {
			name := strings.TrimSpace(s.ChildText(`div:first-child`))
			info := strings.TrimSpace(s.ChildText(`div:nth-child(2)`))
			// spec := records.SpecBase{Spec: make(map[string]string)}
			specs.Spec[name] = info
		})

		product.Detail = &details
		product.Spec = &specs
		product.TimeBase = time
		productDomain := product.ToV1Domain()
		pID, _ := dbProduct.UpsertProduct(&productDomain)

		getReviewStat(prodId, &reviews)
		var avgRating records.AvgRatings
		avgRating.ProductID = pID
		avgRating.Star5 = reviews.Star5
		avgRating.Star4 = reviews.Star4
		avgRating.Star3 = reviews.Star3
		avgRating.Star2 = reviews.Star2
		avgRating.Star1 = reviews.Star1
		avgRating.TimeBase = time

		if len(reviews.BatchedResults.Q0.Results) > 0 {
			avgRating.Overall = reviews.BatchedResults.Q0.Results[0].ReviewStatistics.AverageOverallRating
			avgRating.Overall = float32(math.Round(float64(avgRating.Overall)*10) / 10)
			avgRating.Value = reviews.BatchedResults.Q0.Results[0].ReviewStatistics.SecondaryRatingsAverages.Value.AverageRating
			avgRating.Value = float32(math.Round(float64(avgRating.Value)*10) / 10)
			avgRating.Quality = reviews.BatchedResults.Q0.Results[0].ReviewStatistics.SecondaryRatingsAverages.Quality.AverageRating
			avgRating.Quality = float32(math.Round(float64(avgRating.Quality)*10) / 10)
			avgRating.TotalReviews = reviews.BatchedResults.Q0.Results[0].ReviewStatistics.TotalReviewCount
		} else {
			avgRating.Overall = 0
			avgRating.Value = 0
			avgRating.Quality = 0
			avgRating.TotalReviews = 0
		}

		avgRatingDomain := avgRating.ToV1Domain()
		dbAvgRate.UpsertAvgRating(&avgRatingDomain)

		discHist.ProductID = pID
		discHist.Price = p.Price
		discHist.SaveAmount = p.InstSave
		discHist.SavePercent = float32(math.Round(float64(p.InstSave/p.InWarehouse)*100) / 100)
		discHist.Duration = p.Duration
		discHist.StartedAt = p.StartAt
		discHist.EndedAt = p.EndedAt
		discHist.TimeBase = time

		discHistDomin := discHist.ToV1Domain()
		dbDiscHist.SaveDiscHistory(&discHistDomin)
	})

	d.OnScraped(func(r *colly.Response) {
		logger.Debug("d scraper done: "+r.Request.URL.String(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryScraper})
	})

	// visiting the first page
	c.Visit(pageToScrape)
	// wait for Colly to visit all pages
	c.Wait()

	for _, page := range pagesToScrape {
		logger.Debug("visiting page: "+page, logrus.Fields{})
		// initialize the chan
		done = make(chan struct{})
		d.Visit(page)
		d.Wait()
	}
}
