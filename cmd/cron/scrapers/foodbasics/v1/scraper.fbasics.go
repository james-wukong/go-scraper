package v1

import (
	"math"
	"net/http"
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

func FBasicsScraper(conn *sqlx.DB) {
	var visited int = 0
	var err error
	var pages int
	var pagesToScrape []string
	var prices = make(map[string]scrapers.FBPrices)
	var link string = "https://www.foodbasics.ca/search-page-"
	var apiUrl, nextUrl string
	var params = map[string]string{
		"sortOrder":     "popularity",
		"filter":        ":popularity:deal:Flyer+&+Deals/:deal:Flyer+&+Deals:deal:FLYER_DEAL",
		"fromEcomFlyer": "true",
	}
	var done = make(chan struct{}, 1)
	db := V1Postgres.NewCategoryRepo(conn)
	dbProduct := V1Postgres.NewProductRepo(conn)
	// dbAvgRate := V1Postgres.NewAvgRatingRepo(conn)
	dbDiscHist := V1Postgres.NewDiscHistoryRepo(conn)

	// start from page one
	apiUrl = newFBLink(link, 1, params)

	// creating a new Colly instance
	c := colly.NewCollector(
		// colly.Debugger(&debug.LogDebugger{}),
		// colly.AllowURLRevisit(),
		// turning on the asynchronous request mode in Colly
		colly.Async(true),
	)

	// setting a valid User-Agent header
	c.UserAgent = scrapers.FBasicsUserAgent
	cookies := []*http.Cookie{}
	c.SetCookies(apiUrl, cookies)
	// set delay to prevent from exceeding the rate limit
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*foodbasics.*",
		Parallelism: 2,
		Delay:       1 * time.Second,
		// RandomDelay: 1 * time.Second,
	})
	// scraper for details
	d := c.Clone()
	// Set limits
	// TODO collect total pages to collect items
	// catch the promotion container from coupon link
	c.OnHTML(`div.product-page-nav-standard > div.ppn--pagination > a:nth-last-child(2)`, func(e *colly.HTMLElement) {
		if pages == 0 {
			pages, err = strconv.Atoi(e.Text)
			if err != nil {
				logger.Debug("error", logrus.Fields{"error": err})
			}
			// Signal that pages is done
			close(done)
		}
	})
	// collect basic info from list page
	c.OnHTML(`div.product-page-filter > div[data-list-name="flyerEcomCatalogFilter"] > div.searchOnlineResults`, func(e *colly.HTMLElement) {
		e.ForEach(`div.tile-product`, func(_ int, p *colly.HTMLElement) {
			price := scrapers.Prices{EcoFee: 0.0}
			price.InWarehouse = parsePrice(p.ChildText(`div.pt__content div.content__pricing div.pricing__sale-price`))
			price.Price = parsePrice(p.ChildText(`div.pt__content div.content__pricing div.pricing__before-price`))
			price.InstSave = float32(math.Round((float64(price.Price-price.InWarehouse))*100) / 100)
			// get promotion info
			// get basic price info
			valid := scrapers.Valid{StartAt: time.Now()}
			to := strings.Replace(p.ChildText(`div.pt__content div.content__pricing div.pricing__until-date`), "Valid until ", "", 1)
			valid.EndedAt = parseDate(to)
			valid.Duration = rmWhiteSpaces(to)
			price.Valid = &valid
			// get category info
			fbPrice := scrapers.FBPrices{Prices: price}
			fbPrice.CategoryUrl = p.Attr(`data-category-url`)

			// append link to list
			prices[p.Attr(`data-product-code`)] = fbPrice
			pagesToScrape = append(pagesToScrape, p.ChildAttr(`div.pt__content div.content__head a.product-details-link`, "href"))
		})
	})
	c.OnError(func(r *colly.Response, err error) {
		logger.Debug("Request list URL failed with response", logrus.Fields{"r": string(r.Body), "err": err.Error()})
	})

	// visiting the first page
	c.Visit(apiUrl)
	// block until pages is done
	<-done
	// TODO collect total items links to visit
	for page := 2; page <= pages; page++ {
		nextUrl = newFBLink(link, page, params)
		c.Visit(nextUrl)
	}
	// wait for Colly to visit all pages
	c.Wait()

	// TODO collect item details from items
	// collect detail info from detail page
	d.OnHTML(`div#content-temp`, func(e *colly.HTMLElement) {
		var cId int
		// get product sku
		sku := e.ChildAttr(`div.product-info`, "data-product-code")
		// get category info
		if _, ok := prices[sku]; !ok {
			return
		}

		if ctg, err := db.GetByURLPlatform(prices[sku].CategoryUrl, uint(constants.FOODBASICS)); err != nil {
			logger.Debug("category error caught: ", logrus.Fields{"err": err})
			cId = scrapers.FBasicsDefaultCategory
		} else {
			cId = ctg.Id
		}
		// get product info
		details := records.DetailBase{Detail: make(map[string][]string)}
		specs := records.SpecBase{Spec: make(map[string]string)}
		time := records.TimeBase{
			CreatedAt: time.Now(),
		}
		product := records.Products{Detail: &details, Spec: &specs, TimeBase: time}
		product.Name = rmNewlines(e.ChildText(`div.product-info > div.pi--second-col > div.pi--name > h1.pi--title`))
		product.CategoryId = cId
		product.ProdId = sku
		product.Sku = sku
		product.Brand = e.ChildAttr(`div.product-info`, "data-product-brand")
		// product.Model = e.ChildText(`div#product-details div#product-body-model-number span`)
		product.Price = prices[sku].Price
		product.Source = int(constants.FOODBASICS)
		product.UrlLink = e.Request.URL.String()
		product.ImageLink = e.ChildAttr(`div.pi--first-col picture#main-img img#mob-img`, "src")
		productDomain := product.ToV1Domain()
		pID, _ := dbProduct.UpsertProduct(&productDomain)

		// get discount history
		var discHist records.DiscountHistories
		discHist.ProductID = pID
		discHist.Price = prices[sku].Price
		discHist.SaveAmount = prices[sku].InstSave
		discHist.SavePercent = float32(math.Round(float64(prices[sku].InstSave/prices[sku].InWarehouse)*100) / 100)
		discHist.Duration = prices[sku].Duration
		discHist.StartedAt = prices[sku].StartAt
		discHist.EndedAt = prices[sku].EndedAt
		discHist.TimeBase = time

		discHistDomin := discHist.ToV1Domain()
		dbDiscHist.SaveDiscHistory(&discHistDomin)
	})

	d.OnError(func(r *colly.Response, err error) {
		logger.Debug("Request detail URL failed with response", logrus.Fields{"r": string(r.Body), "err": err.Error()})
	})
	d.OnScraped(func(r *colly.Response) {
		visited += 1
		logger.Debug("page done", logrus.Fields{"idx": visited, "page": r.Request.URL.String()})
	})
	logger.Debug("total page to visit", logrus.Fields{"len": len(pagesToScrape)})
	for _, item := range pagesToScrape {
		if link, err := joinUrlPath(scrapers.FBasicsBaseUrl, item); err != nil {
			logger.Debug("error joining paths", logrus.Fields{"err": err})
		} else {
			d.Visit(link)
		}
	}
	d.Wait()

}
