package v1

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/james-wukong/go-app/cmd/cron/scrapers"
	"github.com/james-wukong/go-app/pkg/logger"
	"github.com/sirupsen/logrus"
)

type reviewResp struct {
	BatchedResults BatchedResults `json:"BatchedResults"`
	*Ratings
}
type BatchedResults struct {
	Q0 Q0 `json:"Q0"`
}

type Q0 struct {
	Results []Results `json:"Results"`
}

type Results struct {
	ReviewStatistics ReviewStatistics `json:"ReviewStatistics"`
}

type ReviewStatistics struct {
	NotRecommendedCount      int                      `json:"NotRecommendedCount"`
	HelpfulVoteCount         int                      `json:"HelpfulVoteCount"`
	NotHelpfulVoteCount      int                      `json:"NotHelpfulVoteCount"`
	RecommendedCount         int                      `json:"RecommendedCount"`
	TotalReviewCount         int                      `json:"TotalReviewCount"`
	RatingsOnlyReviewCount   int                      `json:"RatingsOnlyReviewCount"`
	FeaturedReviewCount      int                      `json:"FeaturedReviewCount"`
	AverageOverallRating     float32                  `json:"AverageOverallRating"`
	SecondaryRatingsAverages SecondaryRatingsAverages `json:"SecondaryRatingsAverages"`
	RatingDistribution       []RatingDist             `json:"RatingDistribution"`
}

type SRADetail struct {
	Id            string  `json:"Id"`
	AverageRating float32 `json:"AverageRating"`
	ValueRange    int     `json:"ValueRange"`
}

type SecondaryRatingsAverages struct {
	Quality SRADetail `json:"Quality"`
	Value   SRADetail `json:"Value"`
}

type RatingDist struct {
	RatingValue int `json:"RatingValue"`
	Count       int `json:"Count"`
}

type Ratings struct {
	Star5 int
	Star4 int
	Star3 int
	Star2 int
	Star1 int
}

// Helper function to parse the price from a string
func parsePrice(priceStr string) (float64, error) {
	parts := strings.Split(priceStr, "$")
	priceStr = strings.Split(parts[1], " ")[0]
	return strconv.ParseFloat(priceStr, 16)
}

func parseDate(rawData string) (t time.Time) {
	var err error
	dt := strings.TrimSpace(rawData)
	if t, err = time.Parse(scrapers.CostcoTimeLayoutISO, dt); err != nil {
		logger.Debug("error parsing date", logrus.Fields{"error": err})
	}
	return
}

func getIDFromUrl(url string) (id string) {
	pattern := `\.(\d{5,})\.html`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		id = match[1]
	}
	return
}

func getNumFromString(input string) string {
	pattern := `\d+`
	re := regexp.MustCompile(pattern)
	match := re.FindString(input)

	return match
}

func getPathFromURL(url string) (name string) {
	pattern := `\/([\w-]{5,})\.htm[l]?`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(url)
	if len(match) > 0 {
		name = match[0]
	}
	return
}

func getReviewStat(id string, result *reviewResp) {
	// make GET request to API to get user by ID
	apiUrl := scrapers.CostcoApiUrl
	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		logger.Debug("error request", logrus.Fields{"error": err})
	}

	// ppending to existing query
	q := request.URL.Query()
	q.Add("passkey", "l7o783yf16tmpcr2d9dwkm783")
	q.Add("apiversion", "5.5")
	q.Add("displaycode", "20040_1_0-en_ca")
	q.Add("resource.q0", "products")
	q.Add("filter.q0", "id:eq:"+id)
	q.Add("stats.q0", "reviews")
	q.Add("filteredstats.q0", "reviews")
	q.Add("filter_reviews.q0", "contentlocale:eq:en_CA,en_US,fr_CA")
	q.Add("filter_reviewcomments.q0", "contentlocale:eq:en_CA,en_US,fr_CA")
	// q.Add("callback", "BV._internal.dataHandler0")
	request.URL.RawQuery = q.Encode()
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	// logger.Debug("api request uri", logrus.Fields{"uri": request.URL.String()})

	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		logger.Debug("error response", logrus.Fields{"error": err})
	}
	// clean up memory after execution
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		logger.Debug("error decoding", logrus.Fields{"error": err})
	}
	if len(result.BatchedResults.Q0.Results) > 0 {
		result.Ratings = &Ratings{}
		for _, rating := range result.BatchedResults.Q0.Results[0].ReviewStatistics.RatingDistribution {
			switch rating.RatingValue {
			case 5:
				result.Ratings.Star5 = rating.Count
			case 4:
				result.Ratings.Star4 = rating.Count
			case 3:
				result.Ratings.Star3 = rating.Count
			case 2:
				result.Ratings.Star2 = rating.Count
			case 1:
				result.Ratings.Star1 = rating.Count
			}
		}
	} else {
		result.Ratings = &Ratings{Star5: 0, Star4: 0, Star3: 0, Star2: 0, Star1: 0}
	}
}
