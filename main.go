package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var demoStrings = []string{
	"..  Opinion: Massive floods have swallowed up the Vermont town I love",
	"..  $20 gift cards and $1 books: GOP primary candidates test novel ways to raise money as they scramble for a spot on next month’s debate stage",
	"..  Russian commander killed while jogging may have been tracked on Strava app",
	"..  Opinion: What Sen. Tommy Tuberville waffling on White nationalism really means",
	"..  Justice Department takes unusual step to try to protect Trump from testifying in lawsuit over FBI firing",
	"..  WAF Awards 2023: World’s best new buildings unveiled",
	"..  Home Run Derby: Vladimir Guerrero Jr. follows in father’s footsteps to win",
}

func main() {
	// Join all the demoStrings into a single string
	joinedString := strings.Join(demoStrings, " ")

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
			// Clear the terminal
			cmd := exec.Command("clear") // For Windows, replace "clear" with "cls"
			cmd.Stdout = os.Stdout
			cmd.Run()

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
			time.Sleep(100 * time.Millisecond)
		}

		// Sleep for a longer duration before scrolling again
		time.Sleep(500 * time.Millisecond)
	}
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
