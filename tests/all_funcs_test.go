package tests

// "log"
// "testing"

// core_functions "github.com/Vanodium/pricetracker/internal/services"
// notifications "github.com/Vanodium/pricetracker/internal/transport"

// func TestEmailNotification(t *testing.T) {
// 	notifications.EmailNotification("", "kverhuka.com")
// 	t.Log("Test email sent!")
// }

// func DownloadHtml(link string) string {
// 	response, err := http.Get(link)
// 	if err != nil {
// 		log.Fatalf("Error fetching the URL: %v", err)
// 	}
// 	defer response.Body.Close() // Ensure the body is closed after we're done

// 	body, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		log.Fatalf("Error reading the response body: %v", err)
// 	}
// 	return string(body)
// }

// func ExtractPrice(link string, path string) string {
// 	htmlContent := DownloadHtml(link)
// 	log.Println(htmlContent)

// 	doc, _ := htmlquery.Parse(strings.NewReader(htmlContent))

// 	// Query for the <h1> element
// 	h1 := htmlquery.FindOne(doc, "/html/head/title")
// 	if h1 != nil {
// 		fmt.Println("title Text:", htmlquery.InnerText(h1))
// 	}

// 	fmt.Println("")
// 	return ""
// }
