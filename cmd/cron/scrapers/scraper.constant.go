package scrapers

const (
	CostcoTimeLayoutISO  string = "2006-01-2" //"January 2, 2006"
	FBasicsTimeLayoutISO string = "Jan 2, 2006"
)

const (
	CostcoBaseUrl         string = "https://www.costco.ca/"
	CostcoApiUrl          string = "https://api.bazaarvoice.com/data/batch.json"
	CostcoUserAgent       string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
	CostcoCookiesName     string = "akavpau_zezxapz5yfca"
	CostcoCookiesValue    string = "1719897913~id=66f6a1ded026631f3dda6119a13f7951"
	CostcoCookiesBmsvKey  string = "bm_sv"
	CostcoCookiesBmsvVal  string = "77B82180A25CFBECC33B31BAA974F87B~YAAQBJYqFz7hBnGQAQAA7vXdcRhW9O6VPR4IHAXX5pK7/nrbCibqk5NUXJqW0S34eK55l8tycOx5HDHGVnegEIJV7WoKfcXXu0mTRYOOSI/RUb3A0cCEpbyjINS8QfQM3Jgiah37+rnCHIiybE1E1euSKoEX1BIZ9gOd9qcLrQ5rzFEPhg8B8Sd+SGbktTgomtFGcUOj0SSBtfuXn69Z83irRHTj89uJ1sAw057aVnH6HdlBxSmiUcgw+qt7Bpa9~1"
	CostcoCookiesBmszKey  string = "bm_sz"
	CostcoCookiesBmszVal  string = "5B843271FF3466AF5D350CDC1FFE44F7~YAAQBJYqFz/hBnGQAQAA7vXdcRgi5fOCCSPWyyMDcKgyt/YPwP5lMaSZNutWQZAIriDjZeI1q6gR51PRDNVhMT32v8tPMO8Cq1Ti7Oi+mburIdKGSw8m3Od3GSCjEPB1xa6LAp36ldwNgekfu0es7yXHVMj/LJtzU4igyiQ7mrGOj27LNSY7ujXnymugxwCiob6Flj3ogNy5rmI7xwmyaM/wiwmvIBCHJtZz8jPt6idDzr8GLTrfdY8ND/JyRThYR25EpKGf+BXowibtJsZJLg8LuvSiYUXsj/6QiJ8DDhPXVN7iWhoLNOfCAO1us24N8NrynYZrIuV7dmH4p8gI+751HPih+yq7uEWosd+RrqHnHS3j+ZZmvSq5mwlcx0Gen8dQja6y54VBdg6JT5OvCQSgw5HtE7x/wZVlB9HoLtzUpgt1YxS6eq+Wq1xoajjp20JoHiy8bBEDW1LCYI/hLKatqQQOj+WVzkFSFQ==~3490617~3162416"
	CostcoDefaultCategory int    = 873
)

const (
	FBasicsBaseUrl         string = "https://www.foodbasics.ca/"
	FBasicsApiUrl          string = "https://www.foodbasics.ca/aisle-ajax"
	FBasicsUserAgent       string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
	FBasicsDefaultCategory int    = 1484
	FBasicsCategoryNum     int    = 50
)
