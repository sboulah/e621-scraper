package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	Init()
}

// Init
func Init() {
	//

	// Create temp directories
	os.Mkdir("./tmp", 0755)

GetInput:
	// New scanner
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Please enter the query tags, beginning after 'tags=': ")

	regexp.Compile(`(?<=\+|\=)[a-z_\%289\(\)]+`)

	// Scan
	scanner.Scan()
	input := scanner.Text()
	tags := input

CheckInput:
	// Scan
	fmt.Printf("You entered: %v, is that correct? Y or N: ", input)
	scanner.Scan()

	// Lowercase
	input = scanner.Text()
	input = strings.ToLower(input)

	// Switch
	switch input {
	case "y", "yes":
		Scrape(tags)
	case "n", "no":
		goto GetInput
	default:
		fmt.Println("I did not understand that!")
		goto CheckInput

	}

}

// Scrape E621
func Scrape(tags string) {
	//

	// Start timer, this will be used to calculate the speed and time taken
	// var mainStopwatch time.Time = time.Now()

	// Page number, how many pages of posts we have found
	var totalPages int = 0

	// // Empty slice
	// var RawPosts []RawPost = []RawPost{}
	// var MetaPosts []MetaPost = []MetaPost{}
	// var CleanedPosts []CleanedPost = []CleanedPost{}

	// var test []owo.RawPost = []owo.RawPost{}

	// Load E621
	LoadE621(tags, &totalPages)

}

// Load E621
func LoadE621(tags string, totalPages *int) {
	//

	// Start
	fmt.Println("Starting...")
	var start time.Time = time.Now()

	// Search URL
	var fullSearchURL string = fmt.Sprintf("https://e621.net/posts?page=%d&tags=%s", 1, tags)

	// Get the page HTML
	response, error := http.Get(fullSearchURL)
	if error != nil {
		log.Fatal(error)
	}

	fmt.Println("Loaded E621")

	// Close the response
	defer response.Body.Close()

	// Make sure the status code is 200
	if response.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s", response.StatusCode, response.Status)
	}

	// Parse the HTML so that goquery can read it
	document, error := goquery.NewDocumentFromReader(response.Body)
	if error != nil {
		log.Fatal(error)
	}

	// Select the paginator
	var paginator *goquery.Selection = document.Find("li.numbered-page")

	// Get the last element of the paginator and convert it to an int
	*totalPages, _ = strconv.Atoi(paginator.Last().Text())

	// Stop
	var end float64 = time.Since(start).Seconds()
	if *totalPages == 0 {
		*totalPages = 1
	}

	// Log
	fmt.Printf("Identified %d page(s) worth of posts in %F seconds!\n", *totalPages, end)
	if *totalPages == 750 {
		fmt.Println("750 is the maximum number of pages E621 allows,\nThe bot will only scrape the first 750 pages.")
	}

}

// // Scrape URLs
// func GetRawPosts(ttlPgNum int, RawPosts *[]RawPost, tags string) {
// 	//

// 	// Loop through all pages
// 	for curPgNum := 1; curPgNum <= ttlPgNum; curPgNum++ {
// 		//

// 		// Set URL to current page number
// 		var fullURL string = fmt.Sprintf("https://e621.net/posts?page=%d&tags=%s", curPgNum, tags)

// 		// Print current page number and URL
// 		fmt.Printf("Current page number: %d, URL: %s\n", curPgNum, fullURL)

// 		// New collector
// 		var c *colly.Collector = colly.NewCollector()

// 		// On HTML
// 		c.OnHTML("article a", func(h *colly.HTMLElement) {
// 			//

// 			// Get URL
// 			var postURL string = h.Attr("href")

// 			// Format URL
// 			var formattedURL string = fmt.Sprintf("https://e621.net%s", strings.Split(postURL, "?q")[0])

// 			// Create struct
// 			post := RawPost{
// 				URL: formattedURL,
// 			}

// 			// Append to slice
// 			*RawPosts = append(*RawPosts, post)

// 			// Print
// 			fmt.Printf("Scraped Post %d: %s\n", len(*RawPosts), formattedURL)
// 		})

// 		// Visit URL
// 		c.Visit(fullURL)
// 	}
// }

// // Save to JSON
// func WriteRawPost(RawPosts *[]RawPost) {
// 	//

// 	// Indent JSON
// 	file, err := json.MarshalIndent(RawPosts, "", " ")
// 	if err != nil {
// 		fmt.Println("Could not create JSON file")
// 		return
// 	}

// 	// Log total posts found
// 	fmt.Printf("Total Posts Found: %d\n", len(*RawPosts))

// 	// Write to file
// 	_ = ioutil.WriteFile("./tmp/rawposts.json", file, 0644)
// }

// // Open JSON
// func OpenRawPosts(RawPosts *[]RawPost) {
// 	//

// 	// Clear slice
// 	*RawPosts = nil

// 	// Open JSON File
// 	data, err := ioutil.ReadFile("./tmp/rawposts.json")
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// Un marshall it
// 	err = json.Unmarshal(data, RawPosts)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// // Scrape Metadata
// func GetPostData(RawPosts *[]RawPost, MetaPosts *[]MetaPost) {
// 	//

// 	// Loop
// 	for i, ii := range *RawPosts {
// 		//

// 		// Current post
// 		fmt.Printf("Processing post %d of %d\n", i+1, len(*RawPosts))

// 		// New collector
// 		var c *colly.Collector = colly.NewCollector()

// 		// On HTML
// 		c.OnHTML("#image-container", func(h *colly.HTMLElement) {
// 			//

// 			// Empty map
// 			var rawJSON MetaPost

// 			// Get attribute
// 			in := []byte(h.Attr("data-post"))

// 			// Un marshall JSON
// 			err := json.Unmarshal(in, &rawJSON)
// 			if err != nil {
// 				panic(err)
// 			}

// 			// Append JSON metadata
// 			*MetaPosts = append(*MetaPosts, rawJSON)

// 		})

// 		// Visit URL
// 		c.Visit(ii.URL)
// 	}
// }

// // Save to JSON
// func WritePostData(MetaPosts *[]MetaPost) {
// 	//

// 	// Indent JSON
// 	file, err := json.MarshalIndent(*MetaPosts, "", " ")
// 	if err != nil {
// 		fmt.Println("Could not create JSON file")
// 		return
// 	}

// 	// Write to file
// 	_ = ioutil.WriteFile("./tmp/rawmetaposts.json", file, 0644)
// }

// // Open JSON
// func OpenPostData(MetaPosts *[]MetaPost) {
// 	//

// 	// Clear slice
// 	*MetaPosts = nil

// 	// Open JSON File
// 	data, err := ioutil.ReadFile("./tmp/rawmetaposts.json")
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// Un marshall it
// 	err = json.Unmarshal(data, MetaPosts)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// // Extract useful data
// func ExtractData(MetaPosts *[]MetaPost, CleanedPosts *[]CleanedPost) {
// 	//

// 	// Loop
// 	for _, i := range *MetaPosts {
// 		//

// 		// Create post
// 		post := CleanedPost{
// 			POSTURL:  fmt.Sprintf("https://e621.net/posts/%d", i.ID),
// 			POSTID:   int(i.ID),
// 			MEDIAURL: i.File.URL,
// 			SOURCES:  i.Sources,
// 			MD5:      i.File.Md5,
// 		}

// 		// Append to slice
// 		*CleanedPosts = append(*CleanedPosts, post)
// 	}
// }

// // Save to JSON
// func WriteCleanedPosts(CleanedPosts *[]CleanedPost) {
// 	//

// 	// Indent JSON
// 	file, err := json.MarshalIndent(*CleanedPosts, "", " ")
// 	if err != nil {
// 		fmt.Println("Could not create JSON file")
// 		return
// 	}

// 	// Write to file
// 	_ = ioutil.WriteFile("./tmp/cleanedposts.json", file, 0644)
// }

// // Open JSON
// func OpenCleanedData(CleanedPosts *[]CleanedPost) {
// 	//

// 	// Clear slice
// 	*CleanedPosts = nil

// 	// Open JSON File
// 	data, err := ioutil.ReadFile("./tmp/cleanedposts.json")
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// Un marshall it
// 	err = json.Unmarshal(data, CleanedPosts)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// // Download Media
// func DownloadMedia(CleanedPosts *[]CleanedPost) {
// 	//

// 	// Create Directory
// 	os.Mkdir("./sauce", 0755)

// 	// Loop
// 	for i, ii := range *CleanedPosts {
// 		//

// 		// Build fileName from fullPath
// 		fileURL, err := url.Parse(ii.MEDIAURL)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		// Set path
// 		segments := strings.Split(fileURL.Path, "/")
// 		fileName := segments[len(segments)-1]

// 		// Create blank file
// 		file, err := os.Create("./sauce/" + fileName)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		// Put content on file
// 		resp, err := http.Get(ii.MEDIAURL)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		defer resp.Body.Close()

// 		io.Copy(file, resp.Body)
// 		defer file.Close()

// 		fmt.Printf("Successfully Downloaded: %s\nNumber %d of %d\n", fileName, i+1, len(*CleanedPosts))

// 	}
// }

// // Cleanup
// func Cleanup(start *time.Time) {
// 	//

// 	// End time
// 	end := time.Since(*start).Seconds()

// 	// Print
// 	fmt.Printf("Completed Scraping, Elapsed Time: %v seconds!\n", end)
// 	fmt.Printf("Press enter to exit! ")

// 	// Scanner
// 	scanner := bufio.NewScanner(os.Stdin)
// 	scanner.Scan()

// 	fmt.Println("Deleting temp files...")
// 	os.Remove("./tmp")
// 	fmt.Println("Deleted, goodbye!")
// }
