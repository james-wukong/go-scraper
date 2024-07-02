package v1

import (
	"math"
	"net/http"
	"regexp"
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

func CostcoScraper(conn *sqlx.DB) {
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
	// set delay to prevent from exceeding the rate limit
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*costco.*",
		Parallelism: 1,
		Delay:       2 * time.Second,
		// RandomDelay: 1 * time.Second,
	})
	// setting a valid User-Agent header
	c.UserAgent = scrapers.CostcoUserAgent
	cookies := []*http.Cookie{
		{
			Name:  scrapers.CostcoCookiesName,
			Value: scrapers.CostcoCookiesValue,
		},
		// {
		// 	Name:  scrapers.CostcoCookiesBmsvKey,
		// 	Value: scrapers.CostcoCookiesBmsvVal,
		// },
		// {
		// 	Name:  scrapers.CostcoCookiesBmszKey,
		// 	Value: scrapers.CostcoCookiesBmszVal,
		// },
		{
			Name:  "Domain",
			Value: "www.costco.ca",
		},
		{
			Name:  "Domain",
			Value: ".costco.ca",
		},
		// {
		// 	Name:  "BCO",
		// 	Value: "pm1",
		// },
		{
			Name:  "Path",
			Value: "/",
		},
		{
			Name:  "SameSite",
			Value: "None",
		},
	}
	c.SetCookies(pageToScrape, cookies)
	// scraper for details
	d := c.Clone()
	var categoryId int
	var prodId string
	// Channel to signal completion of the first callback
	var done chan struct{}
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
				valid := scrapers.Valid{}
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
				logger.Debug("empty link1: ", logrus.Fields{"link": link})
			}
		})
	})

	// catch the catelog search result container
	c.OnHTML(`div.product-list`, func(e *colly.HTMLElement) {
		// discovering a new page
		e.ForEach(`div.product`, func(_ int, p *colly.HTMLElement) {
			link := p.ChildAttr(`span.description a`, "href")

			if strings.Contains(p.Request.URL.String(), searchPattern) {
				if strings.Contains(link, pageDomain) {
					// if this is a product detail link
					pagesToScrape = append(pagesToScrape, link)
				}
			} else {
				logger.Debug("empty link2: ", logrus.Fields{"link2": link})
			}
		})
	})
	c.OnError(func(r *colly.Response, err error) {
		logger.Debug("Request list URL failed with response", logrus.Fields{"err": err.Error(), "r.Body": string(r.Body)})
		// var retries int = 0
		// // Attempt to retry the request
		// for retries < 20 {
		// 	retries += 1
		// 	err = r.Request.Retry()
		// 	if err != nil {
		// 		logger.Debug("retry failed", logrus.Fields{"err": err, "retries": retries})
		// 	} else {
		// 		break
		// 	}
		// }
	})
	// OnResponse callback
	c.OnResponse(func(r *colly.Response) {
		logger.Debug("Request list URL response code", logrus.Fields{"code": r.StatusCode})
	})
	c.OnScraped(func(r *colly.Response) {

	})

	// scrapping the catgories
	d.OnHTML(`ul#crumbs_ul > div > ul > li:last-child a[itemprop="item"]`, func(e *colly.HTMLElement) {
		pathName := getPathFromURL(e.Attr(`href`))
		// categoryD, err := db.GetByNamePlatform(category, uint(constants.COSTCO))
		categoryD, err := db.GetByURLPlatform(pathName, uint(constants.COSTCO))
		categoryId = categoryD.Id
		if err != nil {
			logger.Debug("category error caught: ", logrus.Fields{"category": pathName, "err": err})
			categoryId = scrapers.CostcoDefaultCategory
		}
		close(done) // Signal that the first callback is done
	})

	// scraping promotion details
	d.OnHTML(`div[itemtype="https://schema.org/Product"]`, func(e *colly.HTMLElement) {
		// logger.Debug("receiving signal chan: ", logrus.Fields{})
		<-done // Wait for the signal from the first callback
		time := records.TimeBase{
			CreatedAt: time.Now(),
		}
		// get product info
		product := records.Products{TimeBase: time}
		pName := e.ChildText(`div#product-details div.product-h1-container-v2 h1[itemprop="name"]`)
		product.Name = strings.Split(pName, ", ")[0]
		product.CategoryId = categoryId
		product.ProdId = prodId
		itemSku := getIDFromUrl(e.Request.URL.String())
		prodId = getIDFromUrl(e.Request.URL.String())
		// product.Sku = getIDFromUrl(e.Request.URL.String())
		product.Sku = getNumFromString(e.ChildText(`div#product-details div#product-body-item-number span`))

		// Check if product.Sku exists in the price map
		p, ok := price[itemSku]
		if !ok {
			// If product.Sku does not exist, check if itemSku exists in the price map
			p, ok = price[product.Sku]
			if !ok {
				// If neither exists, skip the rest of the code
				logger.Debug("price & sku not found: ", logrus.Fields{"sku": product.Sku, "price": p})
				return
			}
		}
		product.Model = e.ChildText(`div#product-details div#product-body-model-number span`)
		product.Price = p.InWarehouse
		product.Source = int(constants.COSTCO)
		product.UrlLink = e.Request.URL.String()
		product.ImageLink = e.ChildAttr(`div#left-side-content > div#zoomViewer > div#productImageContainer > div#productImageOverlay > img#initialProductImage`, "src")

		descNode := e.DOM.Find(`div#product-info > div > div > div > div > div#nav-pdp-tab-header-3 > div > div.product-info-description > span#productDescriptions1`)
		description := strings.TrimSpace(descNode.Contents().First().Text())
		// get product details
		details := records.DetailBase{Detail: make(map[string][]string)}
		details.Detail["Description"] = []string{description}
		// TODO add more details into struct
		// next := descNode.Text()
		// get product specification
		specs := records.SpecBase{Spec: make(map[string]string)}
		e.ForEach(`div#product-info > div > div > div > div > div#nav-pdp-tab-header-5 > div > div.product-info-description > div > div`, func(_ int, s *colly.HTMLElement) {
			name := strings.TrimSpace(s.ChildText(`div:first-child`))
			info := strings.TrimSpace(s.ChildText(`div:nth-child(2)`))
			// spec := records.SpecBase{Spec: make(map[string]string)}
			specs.Spec[name] = info
			if strings.ToLower(name) == "brand" {
				product.Brand = info
			}
		})

		product.Detail = &details
		product.Spec = &specs
		product.TimeBase = time
		productDomain := product.ToV1Domain()
		pID, _ := dbProduct.UpsertProduct(&productDomain)

		// get avg ratings
		var reviews reviewResp
		var avgRating records.AvgRatings
		getReviewStat(prodId, &reviews)
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

		// get discout history
		var discHist records.DiscountHistories
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
	d.OnError(func(r *colly.Response, err error) {
		logger.Debug("Request detail URL failed with response", logrus.Fields{"r": string(r.Body), "err": err.Error()})
	})
	// OnResponse callback
	d.OnResponse(func(r *colly.Response) {
		logger.Debug("Request detail URL response code", logrus.Fields{"code": r.StatusCode})
	})
	d.OnScraped(func(r *colly.Response) {

	})

	// visiting the first page
	c.Visit(pageToScrape)
	// wait for Colly to visit all pages
	c.Wait()

	// logger.Debug("all prices: ", logrus.Fields{"price": price})

	for _, page := range pagesToScrape {
		logger.Debug("visiting page: ", logrus.Fields{"page": page})
		// initialize the chan
		done = make(chan struct{})
		d.Visit(page)
		d.Wait()
		// if idx >= 3 {
		// 	break
		// }
	}
}
