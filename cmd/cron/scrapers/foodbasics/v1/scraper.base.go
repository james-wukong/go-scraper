package v1

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/james-wukong/go-app/cmd/cron/scrapers"
	V1Domains "github.com/james-wukong/go-app/internal/business/domains/v1"
	v1 "github.com/james-wukong/go-app/internal/business/domains/v1"
	"github.com/james-wukong/go-app/internal/constants"
	"github.com/james-wukong/go-app/pkg/logger"
	"github.com/sirupsen/logrus"
)

type Category struct {
	Id       string      `json:"id"`
	Name     string      `json:"name"`
	Url      string      `json:"url"`
	Children []*Category `json:"children"`
}

func mapRecord(r Category, pid *int, level uint) *V1Domains.CategoryDomain {
	return &V1Domains.CategoryDomain{
		ParentId: pid,
		Name:     r.Name,
		Level:    level,
		Url:      r.Url,
		Platform: uint(constants.FOODBASICS),
	}
}

func readCategory(db v1.CategoryRepoInterface, r *Category, pid *int, level uint) error {
	var child *Category
	var cid int
	var err error
	// TODO save category
	c := mapRecord(*r, pid, level)
	if cid, err = db.Upsert(c); err != nil {
		logger.Debug("error in readCategory", logrus.Fields{"err": err})
		return err
	} else {
		pid = &cid
	}
	// loop through child
	if r.Children != nil {
		level += 1
		for _, child = range r.Children {
			readCategory(db, child, pid, level)
		}
	}
	return nil
}

func newFBLink(base string, page int, p ...map[string]string) string {
	apiUrl := base + strconv.Itoa(page)
	params := url.Values{}
	for _, v := range p {
		for key, value := range v {
			params.Add(key, value)
		}
	}
	apiUrl += "?" + params.Encode()
	return apiUrl
}

func parseDate(rawData string) time.Time {
	rawData = rmWhiteSpaces(rawData)
	dt := strings.TrimSpace(rawData)
	t, err := time.Parse(scrapers.FBasicsTimeLayoutISO, dt)
	if err != nil {
		logger.Debug("error parsing date", logrus.Fields{"error": err})
	}
	return t
}

func parsePrice(strPrice string) float32 {
	pattern := `\d+.\d{2}`
	re := regexp.MustCompile(pattern)
	output := re.FindString(strPrice)
	value, err := strconv.ParseFloat(output, 32)
	if err != nil {
		logger.Debug("Error converting string to float", logrus.Fields{"err": err})
	}

	return float32(value)
}

func rmWhiteSpaces(input string) string {
	pattern := `[^a-zA-Z0-9, ]`
	re := regexp.MustCompile(pattern)
	output := re.ReplaceAllString(input, "")
	output = strings.TrimSpace(output)
	return output
}

func rmNewlines(input string) string {
	pattern := `[\n\r]`
	re := regexp.MustCompile(pattern)
	output := re.ReplaceAllString(input, "")
	output = strings.TrimSpace(output)
	return output
}

func joinUrlPath(base string, paths ...string) (string, error) {
	if s, err := url.JoinPath(base, paths...); err != nil {
		logger.Debug("Error parsing url", logrus.Fields{"err": err})
		return "", err
	} else {
		return s, nil
	}
}
