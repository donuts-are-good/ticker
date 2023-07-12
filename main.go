package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/mmcdole/gofeed"
)

const feedsFile = ".tickerfeeds.list"

type Item struct {
	PubDate time.Time
	Title   string
	Link    string
}

func main() {
	// Get the path to the feeds file in the XDG home directory
	feedsFilePath, err := getFeedsFilePath()
	if err != nil {
		fmt.Println("Failed to get feeds file path:", err)
		return
	}

	var items []Item
	for _, feedURL := range rssFeeds {
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(feedURL)
		if err != nil {
			fmt.Printf("Failed to fetch RSS feed from %s: %v\n", feedURL, err)
			continue
		}

		for _, item := range feed.Items {
			pubDate, err := parsePublishedTime(item.Published)
			if err != nil {
				fmt.Printf("Failed to parse published time from %s: %v\n", item.Published, err)
				continue
			}

			items = append(items, Item{PubDate: pubDate, Title: item.Title, Link: item.Link})
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].PubDate.Before(items[j].PubDate)
	})

	for _, item := range items {
		headline := fmt.Sprintf("%s \"%s\" %s  ..  ", item.PubDate.Format("2006-01-02 15:04:05"), item.Title, item.Link)
		headlines = append(headlines, headline)
	}

	// Join all the headlines into a single string
	joinedString := strings.Join(headlines, "")

	// Get the terminal width to fit the scrolling text
	terminalWidth, err := getTerminalWidth()
	if err != nil {
		fmt.Println("Failed to get terminal width:", err)
		return
	}

	// Determine the maximum length of the scrolling line
	maxLength := terminalWidth

	// Start scrolling loop
	for {
		// Scroll the joined string
		for i := 0; i <= len(joinedString); i++ {
			// Clear the terminal based on the operating system
			clearTerminal()

			// Build the scrolling line
			var scrollingLine string
			if maxLength-i > 0 {
				scrollingLine = strings.Repeat(" ", maxLength-i) + joinedString[:i]
			} else {
				scrollingLine = joinedString[i-maxLength : i]
			}

			// Print the scrolling line
			fmt.Print(scrollingLine)

			// Sleep for a short duration to control the scrolling speed
			time.Sleep(200 * time.Millisecond)
		}
	}
}

// getFeedsFilePath returns the path to the feeds file in the XDG home directory
func getFeedsFilePath() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return homeDir + "/.config/" + feedsFile, nil
}

// readRSSFeeds reads the RSS feeds from the given file path
func readRSSFeeds(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var feeds []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		feed := strings.TrimSpace(scanner.Text())
		if feed != "" && !strings.HasPrefix(feed, "#") {
			feeds = append(feeds, feed)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}

// getTerminalWidth retrieves the width of the terminal window
func getTerminalWidth() (int, error) {
	cmd := exec.Command("tput", "cols")
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	width := strings.TrimSpace(string(out))
	return strconv.Atoi(width)
}

// clearTerminal clears the terminal based on the operating system
func clearTerminal() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
