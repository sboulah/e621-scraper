package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	Init()
}

//////////////////////////////////////////////////////////////////////////////

type RawPost struct {
	URL string `json:"url"`
}

type CleanedPost struct {
	POSTURL  string        `json:"post_url"`
	POSTID   int           `json:"id"`
	MEDIAURL string        `json:"media_url"`
	SOURCES  []interface{} `json:"sources"`
	MD5      string        `json:"md5"`
}

type MetaPost struct {
	ChangeSeq     int64         `json:"change_seq"`
	CommentCount  int64         `json:"comment_count"`
	CreatedAt     string        `json:"created_at"`
	Description   string        `json:"description"`
	FavCount      int64         `json:"fav_count"`
	File          File          `json:"file"`
	Flags         Flags         `json:"flags"`
	ID            int64         `json:"id"`
	IsFavorited   bool          `json:"is_favorited"`
	LockedTags    []interface{} `json:"locked_tags"`
	Pools         []interface{} `json:"pools"`
	Relationships Relationships `json:"relationships"`
	Sample        Sample        `json:"sample"`
	Score         Score         `json:"score"`
	Sources       []interface{} `json:"sources"`
	Tags          []interface{} `json:"tags"`
	UpdatedAt     string        `json:"updated_at"`
	UploaderID    int64         `json:"uploader_id"`
}

type File struct {
	EXT    string `json:"ext"`
	Height int64  `json:"height"`
	Md5    string `json:"md5"`
	Size   int64  `json:"size"`
	URL    string `json:"url"`
	Width  int64  `json:"width"`
}

type Flags struct {
	Deleted      bool `json:"deleted"`
	Flagged      bool `json:"flagged"`
	HasNotes     bool `json:"has_notes"`
	NoteLocked   bool `json:"note_locked"`
	Pending      bool `json:"pending"`
	RatingLocked bool `json:"rating_locked"`
	StatusLocked bool `json:"status_locked"`
}

type Relationships struct {
	Children          []interface{} `json:"children"`
	HasActiveChildren bool          `json:"has_active_children"`
	HasChildren       bool          `json:"has_children"`
	ParentID          interface{}   `json:"parent_id"`
}

type Sample struct {
	Alternates Alternates `json:"alternates"`
	Has        bool       `json:"has"`
	Height     int64      `json:"height"`
	URL        string     `json:"url"`
	Width      int64      `json:"width"`
}

type Alternates struct {
}

type Score struct {
	Down  int64 `json:"down"`
	Total int64 `json:"total"`
	Up    int64 `json:"up"`
}

//////////////////////////////////////////////////////////////////////////////

// Init
func Init() {
	//

	// Create directories
	os.Mkdir("./tmp", 0755)

GetInput:
	// New scanner
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Please enter the query tags, beginning after 'tags=': ")

	// Scan
	scanner.Scan()
	input := scanner.Text()
	tags := input // krystal+loona_%28vivzmind%29++

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

func Scrape(tags string) {
	//

	start := time.Now()

	// Page number
	var ttlPgNum int = 0

	// Empty slice
	var RawPosts []RawPost = []RawPost{}
	var MetaPosts []MetaPost = []MetaPost{}
	var CleanedPosts []CleanedPost = []CleanedPost{}

	// Load E621
	LoadE621(tags, &ttlPgNum)

	// Scrape URLs
	GetRawPosts(ttlPgNum, &RawPosts, tags)

	// Save to JSON
	WriteRawPost(&RawPosts)

	// Open JSON
	OpenRawPosts(&RawPosts)

	// Scrape Metadata
	GetPostData(&RawPosts, &MetaPosts)

	// Save to JSON
	WritePostData(&MetaPosts)

	// Open JSON
	OpenPostData(&MetaPosts)

	// Extract useful data
	ExtractData(&MetaPosts, &CleanedPosts)

	// Save to JSON
	WriteCleanedPosts(&CleanedPosts)

	// Open JSON
	OpenCleanedData(&CleanedPosts)

	// Download Media
	DownloadMedia(&CleanedPosts)

	// Cleanup
	Cleanup(&start)
}

// Load E621
func LoadE621(tags string, ttlPgNum *int) {
	//

	// Full URL
	var fullURL string = fmt.Sprintf("https://e621.net/posts?page=%d&tags=%s", 1, tags)

	// Create a new collector
	var c *colly.Collector = colly.NewCollector()

	// Get the max number of pages
	c.OnHTML("li.numbered-page a", func(h *colly.HTMLElement) {
		*ttlPgNum, _ = strconv.Atoi(h.Text)
	})

	// Visit the page
	c.Visit(fullURL)

	// Load E621
	fmt.Println("Loaded E621")

	if *ttlPgNum == 0 {
		*ttlPgNum = 1
	}

	// Print total pages found
	fmt.Printf("Total Pages Found: %d\n", *ttlPgNum)
}

// Scrape URLs
func GetRawPosts(ttlPgNum int, RawPosts *[]RawPost, tags string) {
	//

	// Loop through all pages
	for curPgNum := 1; curPgNum <= ttlPgNum; curPgNum++ {
		//

		// Set URL to current page number
		var fullURL string = fmt.Sprintf("https://e621.net/posts?page=%d&tags=%s", curPgNum, tags)

		// Print current page number and URL
		fmt.Printf("Current page number: %d, URL: %s\n", curPgNum, fullURL)

		// New collector
		var c *colly.Collector = colly.NewCollector()

		// On HTML
		c.OnHTML("article a", func(h *colly.HTMLElement) {
			//

			// Get URL
			var postURL string = h.Attr("href")

			// Format URL
			var formattedURL string = fmt.Sprintf("https://e621.net%s", strings.Split(postURL, "?q")[0])

			// Create struct
			post := RawPost{
				URL: formattedURL,
			}

			// Append to slice
			*RawPosts = append(*RawPosts, post)

			// Print
			fmt.Printf("Scraped Post %d: %s\n", len(*RawPosts), formattedURL)
		})

		// Visit URL
		c.Visit(fullURL)
	}
}

// Save to JSON
func WriteRawPost(RawPosts *[]RawPost) {
	//

	// Indent JSON
	file, err := json.MarshalIndent(RawPosts, "", " ")
	if err != nil {
		fmt.Println("Could not create JSON file")
		return
	}

	// Log total posts found
	fmt.Printf("Total Posts Found: %d\n", len(*RawPosts))

	// Write to file
	_ = ioutil.WriteFile("./tmp/rawposts.json", file, 0644)
}

// Open JSON
func OpenRawPosts(RawPosts *[]RawPost) {
	//

	// Clear slice
	*RawPosts = nil

	// Open JSON File
	data, err := ioutil.ReadFile("./tmp/rawposts.json")
	if err != nil {
		fmt.Println(err)
	}

	// Un marshall it
	err = json.Unmarshal(data, RawPosts)
	if err != nil {
		fmt.Println(err)
	}
}

// Scrape Metadata
func GetPostData(RawPosts *[]RawPost, MetaPosts *[]MetaPost) {
	//

	// Loop
	for i, ii := range *RawPosts {
		//

		// Current post
		fmt.Printf("Processing post %d of %d\n", i+1, len(*RawPosts))

		// New collector
		var c *colly.Collector = colly.NewCollector()

		// On HTML
		c.OnHTML("#image-container", func(h *colly.HTMLElement) {
			//

			// Empty map
			var rawJSON MetaPost

			// Get attribute
			in := []byte(h.Attr("data-post"))

			// Un marshall JSON
			err := json.Unmarshal(in, &rawJSON)
			if err != nil {
				panic(err)
			}

			// Append JSON metadata
			*MetaPosts = append(*MetaPosts, rawJSON)

		})

		// Visit URL
		c.Visit(ii.URL)
	}
}

// Save to JSON
func WritePostData(MetaPosts *[]MetaPost) {
	//

	// Indent JSON
	file, err := json.MarshalIndent(*MetaPosts, "", " ")
	if err != nil {
		fmt.Println("Could not create JSON file")
		return
	}

	// Write to file
	_ = ioutil.WriteFile("./tmp/rawmetaposts.json", file, 0644)
}

// Open JSON
func OpenPostData(MetaPosts *[]MetaPost) {
	//

	// Clear slice
	*MetaPosts = nil

	// Open JSON File
	data, err := ioutil.ReadFile("./tmp/rawmetaposts.json")
	if err != nil {
		fmt.Println(err)
	}

	// Un marshall it
	err = json.Unmarshal(data, MetaPosts)
	if err != nil {
		fmt.Println(err)
	}
}

// Extract useful data
func ExtractData(MetaPosts *[]MetaPost, CleanedPosts *[]CleanedPost) {
	//

	// Loop
	for _, i := range *MetaPosts {
		//

		// Create post
		post := CleanedPost{
			POSTURL:  fmt.Sprintf("https://e621.net/posts/%d", i.ID),
			POSTID:   int(i.ID),
			MEDIAURL: i.File.URL,
			SOURCES:  i.Sources,
			MD5:      i.File.Md5,
		}

		// Append to slice
		*CleanedPosts = append(*CleanedPosts, post)
	}
}

// Save to JSON
func WriteCleanedPosts(CleanedPosts *[]CleanedPost) {
	//

	// Indent JSON
	file, err := json.MarshalIndent(*CleanedPosts, "", " ")
	if err != nil {
		fmt.Println("Could not create JSON file")
		return
	}

	// Write to file
	_ = ioutil.WriteFile("./tmp/cleanedposts.json", file, 0644)
}

// Open JSON
func OpenCleanedData(CleanedPosts *[]CleanedPost) {
	//

	// Clear slice
	*CleanedPosts = nil

	// Open JSON File
	data, err := ioutil.ReadFile("./tmp/cleanedposts.json")
	if err != nil {
		fmt.Println(err)
	}

	// Un marshall it
	err = json.Unmarshal(data, CleanedPosts)
	if err != nil {
		fmt.Println(err)
	}
}

// Download Media
func DownloadMedia(CleanedPosts *[]CleanedPost) {
	//

	// Create Directory
	os.Mkdir("./sauce", 0755)

	// Loop
	for i, ii := range *CleanedPosts {
		//

		// Build fileName from fullPath
		fileURL, err := url.Parse(ii.MEDIAURL)
		if err != nil {
			log.Println(err)
		}

		// Set path
		segments := strings.Split(fileURL.Path, "/")
		fileName := segments[len(segments)-1]

		// Create blank file
		file, err := os.Create("./sauce/" + fileName)
		if err != nil {
			log.Println(err)
		}

		// Put content on file
		resp, err := http.Get(ii.MEDIAURL)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()

		io.Copy(file, resp.Body)
		defer file.Close()

		fmt.Printf("Successfully Downloaded: %s\nNumber %d of %d\n", fileName, i+1, len(*CleanedPosts))

	}
}

// Cleanup
func Cleanup(start *time.Time) {
	//

	// End time
	end := time.Since(*start).Seconds()

	// Print
	fmt.Printf("Completed Scraping, Elapsed Time: %v seconds!\n", end)
	fmt.Printf("Press enter to exit! ")

	// Scanner
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	fmt.Println("Deleting temp files...")
	os.Remove("./tmp")
	fmt.Println("Deleted, goodbye!")
}
