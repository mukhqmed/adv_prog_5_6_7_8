package filtering

import (
	"fmt"
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func TestFilteredBarbersHandlerWithSelenium(t *testing.T) {
	// Start Selenium WebDriver for Internet Explorer
	caps := selenium.Capabilities{"browserName": "internet explorer"}
	caps["ie.ensureCleanSession"] = true

	// Connect to Selenium WebDriver server running on port 5555
	wd, err := selenium.NewRemote(caps, "http://localhost:4444/wd/hub")
	if err != nil {
		t.Fatalf("Failed to open session: %s\n", err)
	}
	defer wd.Quit()

	// Open the web page
	if err := wd.Get("http://localhost:8080"); err != nil {
		t.Fatalf("Failed to load page: %s\n", err)
	}

	// Interact with the web page
	searchInput, err := wd.FindElement(selenium.ByID, "status")
	if err != nil {
		t.Fatalf("Failed to find search input: %s\n", err)
	}
	searchInput.SendKeys("Senior")

	experienceInput, err := wd.FindElement(selenium.ByID, "experience")
	if err != nil {
		t.Fatalf("Failed to find experience input: %s\n", err)
	}
	experienceInput.SendKeys("6 лет")

	sortByInput, err := wd.FindElement(selenium.ByID, "sort")
	if err != nil {
		t.Fatalf("Failed to find sort input: %s\n", err)
	}
	sortByInput.SendKeys("name")

	pageInput, err := wd.FindElement(selenium.ByID, "page")
	if err != nil {
		t.Fatalf("Failed to find page input: %s\n", err)
	}
	pageInput.SendKeys("1")

	searchButton, err := wd.FindElement(selenium.ByID, "search-button")
	if err != nil {
		t.Fatalf("Failed to find search button: %s\n", err)
	}
	searchButton.Click()

	time.Sleep(2 * time.Second)

	searchResults, err := wd.FindElements(selenium.ByClassName, "search-result")
	if err != nil {
		t.Fatalf("Failed to find search results: %s\n", err)
	}

	if len(searchResults) == 0 {
		t.Error("No search results found")
	}

	fmt.Println("Test passed: Found search results")
}
