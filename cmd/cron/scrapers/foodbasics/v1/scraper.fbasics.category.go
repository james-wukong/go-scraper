package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/james-wukong/go-app/cmd/cron/scrapers"
	V1Postgres "github.com/james-wukong/go-app/internal/datasources/repositories/postgres/v1"
	"github.com/james-wukong/go-app/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func FBasicsCategoryScraper(conn *sqlx.DB) {
	// make GET request to API to get user by ID
	var request *http.Request
	var response *http.Response
	var result *Category
	var err error
	var q url.Values
	apiUrl := scrapers.FBasicsApiUrl
	client := &http.Client{Timeout: 30 * time.Second}
	db := V1Postgres.NewCategoryRepo(conn)
	// recover panic
	defer func() {
		if r := recover(); r != nil {
			logger.Debug("recovering from panic", logrus.Fields{"panic": r})
		}
	}()
	// loop through different ids
	for i := 1; i <= scrapers.FBasicsCategoryNum; i++ {
		request, err = http.NewRequest(http.MethodGet, apiUrl, nil)
		if err != nil {
			logger.Debug("error request", logrus.Fields{"error": err})
		}

		// appending to existing query
		q = request.URL.Query()

		q.Add("aisleId", fmt.Sprintf("%06d", i))
		// q.Add("callback", "BV._internal.dataHandler0")
		request.URL.RawQuery = q.Encode()
		request.Header.Set("Content-Type", "application/json; charset=utf-8")
		// logger.Debug("api request uri", logrus.Fields{"uri": request.URL.String()})

		response, err = client.Do(request)
		if err != nil {
			logger.Debug("error response", logrus.Fields{"error": err})
		}
		// clean up memory after execution
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
				logger.Debug("error decoding", logrus.Fields{"error": err})
				panic(fmt.Sprintf("error decoding: %s", err.Error()))
			} else if result != nil {
				// TODO save category
				if err = readCategory(db, result, nil, uint(0)); err != nil {
					logger.Debug("error reading category", logrus.Fields{"error": err})
				}
			}
		}
	}
}
