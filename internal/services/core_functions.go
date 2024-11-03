package core_functions

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	db_operations "github.com/Vanodium/pricetracker/internal/db"
	notifications "github.com/Vanodium/pricetracker/internal/transport"

	// "github.com/antchfx/htmlquery"
	"github.com/PuerkitoBio/goquery"
)

func GetCurrentDate() int64 {
	return time.Now().Unix()
}

func CheckPrices() {
	currentDate := GetCurrentDate()
	oldPrices := db_operations.GetOldPrices(currentDate)

	var trackerId, oldPrice, currentPrice string
	var trackerIdInt int64

	var userId int64
	var link, xPath string

	for _, recording := range oldPrices {
		trackerId, oldPrice = recording[0], recording[1]
		trackerIdInt, _ = strconv.ParseInt(trackerId, 10, 64)
		userId, link, xPath = db_operations.GetTrackerById(trackerIdInt)

		currentPrice = ParseDigits(ExtractPrice(link, xPath))
		if currentPrice != oldPrice {
			userEmail := db_operations.GetEmailById(userId)
			notifications.EmailNotification(userEmail, link)
			db_operations.UpdatePrice(trackerIdInt, currentPrice)
		}
		db_operations.UpdatePriceDate(trackerIdInt, currentDate)
	}
}

// func ExtractPrice(link string, path string) string {
// 	doc, err := htmlquery.LoadURL(link)
// 	if err != nil {
// 		panic(err)
// 	}
// 	list, err := htmlquery.QueryAll(doc, path)
// 	if err != nil {
// 		panic(err)
// 	}
// 	if list != nil {
// 		log.Println("Found data by xPath and link")
// 		return htmlquery.InnerText(list[0])
// 	}
// 	log.Println("Did not find price. Captcha?")
// 	return ""
// }

func ExtractPrice(link string, path string) string {
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var options []string
	doc.Find(path).Each(func(i int, s *goquery.Selection) {
		options = append(options, s.Text())
	})
	return options[0]
}

func ParseDigits(block string) string {
	price := regexp.MustCompile("[0-9]+").FindAllString(block, -1)
	return strings.Join(price, "")
}
