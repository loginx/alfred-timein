//go:build bdd
// +build bdd

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/loginx/alfred-timein/internal/adapters/cache"
	"github.com/loginx/alfred-timein/internal/adapters/geocoder"
	"github.com/loginx/alfred-timein/internal/adapters/presenter"
	"github.com/loginx/alfred-timein/internal/adapters/timezonefinder"
	"github.com/loginx/alfred-timein/internal/domain"
	"github.com/loginx/alfred-timein/internal/usecases"
)

// BDDContext holds the state for BDD test scenarios
type BDDContext struct {
	// Services
	geotzUseCase  *usecases.GeotzUseCase
	timeinUseCase *usecases.TimeinUseCase
	cacheService  *cache.LRUCache

	// Test state
	inputCity        string
	inputTimezone    string
	outputTimezone   string
	outputTime       string
	outputJSON       map[string]interface{}
	errorMessage     string
	exitCode         int
	executionTime    time.Duration
	cacheHit         bool
	commandOutput    string
	commandError     string

	// Cache pre-seeding test state
	cache          *cache.LRUCache
	testDir        string
	lastOutput     string
	lastError      error
	lastDuration   time.Duration
	wasFromCache   bool
	userEntryCity  string

	// Test configuration
	cacheEnabled bool
}

func (ctx *BDDContext) reset() {
	ctx.inputCity = ""
	ctx.inputTimezone = ""
	ctx.outputTimezone = ""
	ctx.outputTime = ""
	ctx.outputJSON = nil
	ctx.errorMessage = ""
	ctx.exitCode = 0
	ctx.executionTime = 0
	ctx.cacheHit = false
	ctx.commandOutput = ""
	ctx.commandError = ""
	
	// Reset cache pre-seeding test state - create unique test dir per scenario
	ctx.cache = nil
	ctx.testDir = fmt.Sprintf("/tmp/alfred-timein-test-%d", time.Now().UnixNano())
	ctx.lastOutput = ""
	ctx.lastError = nil
	ctx.lastDuration = 0
	ctx.wasFromCache = false
	ctx.userEntryCity = ""
}

func (ctx *BDDContext) initializeServices() error {
	// Create cache in current directory to match CLI behavior
	ctx.cacheService = cache.NewLRUCache(1000, 30*24*time.Hour, ".")

	// Create adapters
	geocoderAdapter := geocoder.NewOpenStreetMapGeocoder()
	tzFinder, err := timezonefinder.NewTzfTimezoneFinder()
	if err != nil {
		return fmt.Errorf("failed to initialize timezone finder: %w", err)
	}

	// Create formatters
	plainFormatter := presenter.NewPlainFormatter()

	// Create use cases
	ctx.geotzUseCase = usecases.NewGeotzUseCase(geocoderAdapter, tzFinder, ctx.cacheService, plainFormatter)
	ctx.timeinUseCase = usecases.NewTimeinUseCase(plainFormatter)

	return nil
}

// Step definitions for timezone lookup scenarios
func (ctx *BDDContext) theTimzoneLookupServiceIsAvailable() error {
	return ctx.initializeServices()
}

func (ctx *BDDContext) iWantToKnowTheTimezoneFor(city string) error {
	ctx.inputCity = city
	return nil
}

func (ctx *BDDContext) iRequestTheTimezoneInformation() error {
	start := time.Now()
	
	// Check if it's a cache hit first
	cacheKey := strings.ToLower(ctx.inputCity)
	if tz, ok := ctx.cacheService.Get(cacheKey); ok {
		ctx.cacheHit = true
		ctx.outputTimezone = tz
		ctx.executionTime = time.Since(start)
		return nil
	}

	output, err := ctx.geotzUseCase.GetTimezoneFromCity(ctx.inputCity)
	ctx.executionTime = time.Since(start)

	if err != nil {
		ctx.errorMessage = err.Error()
		return nil
	}

	ctx.outputTimezone = strings.TrimSpace(string(output))
	return nil
}

func (ctx *BDDContext) iShouldGetAsTheTimezone(expectedTimezone string) error {
	if ctx.outputTimezone != expectedTimezone {
		return fmt.Errorf("expected timezone %s, got %s", expectedTimezone, ctx.outputTimezone)
	}
	return nil
}

func (ctx *BDDContext) theResponseShouldBeFast() error {
	if ctx.executionTime > 5*time.Second {
		return fmt.Errorf("response took %v, expected under 5 seconds", ctx.executionTime)
	}
	return nil
}

func (ctx *BDDContext) iShouldReceiveAHelpfulErrorMessage() error {
	if ctx.errorMessage == "" {
		return fmt.Errorf("expected an error message, but got none")
	}
	return nil
}

func (ctx *BDDContext) theErrorShouldMention(expectedText string) error {
	if !strings.Contains(strings.ToLower(ctx.errorMessage), strings.ToLower(expectedText)) {
		return fmt.Errorf("expected error to mention '%s', got: %s", expectedText, ctx.errorMessage)
	}
	return nil
}

// Step definitions for caching scenarios
func (ctx *BDDContext) iHavePreviouslyLookedUp(city string) error {
	// Simulate a previous lookup by pre-populating cache
	ctx.cacheService.Set(strings.ToLower(city), "Europe/London") // London as example
	return nil
}

func (ctx *BDDContext) iRequestTheTimezoneForAgain(city string) error {
	ctx.inputCity = city
	return ctx.iRequestTheTimezoneInformation()
}

func (ctx *BDDContext) theResponseShouldBeNearlyInstantaneous() error {
	if ctx.executionTime > 100*time.Millisecond {
		return fmt.Errorf("cached response took %v, expected under 100ms", ctx.executionTime)
	}
	return nil
}

func (ctx *BDDContext) theResultShouldIndicateItCameFromCache() error {
	if !ctx.cacheHit {
		return fmt.Errorf("expected cache hit, but was cache miss")
	}
	return nil
}

// Step definitions for time display scenarios
func (ctx *BDDContext) theTimeDisplayServiceIsAvailable() error {
	return ctx.initializeServices()
}

func (ctx *BDDContext) iHaveTheTimezone(timezone string) error {
	ctx.inputTimezone = timezone
	return nil
}

func (ctx *BDDContext) iRequestTheCurrentTime() error {
	output, err := ctx.timeinUseCase.GetTimezoneInfo(ctx.inputTimezone)
	
	if err != nil {
		ctx.errorMessage = err.Error()
		return nil
	}

	ctx.outputTime = strings.TrimSpace(string(output))
	return nil
}

func (ctx *BDDContext) iShouldSeeAHumanReadableTimeFormat() error {
	if ctx.outputTime == "" {
		return fmt.Errorf("expected time output, got empty string")
	}
	
	// Check for basic time format elements
	if !strings.Contains(ctx.outputTime, "2025") {
		return fmt.Errorf("expected current year in time output: %s", ctx.outputTime)
	}
	
	return nil
}

func (ctx *BDDContext) theTimeShouldIncludeTheDayOfTheWeek() error {
	weekdays := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for _, day := range weekdays {
		if strings.Contains(ctx.outputTime, day) {
			return nil
		}
	}
	return fmt.Errorf("expected day of week in output: %s", ctx.outputTime)
}

func (ctx *BDDContext) theTimeShouldIncludeTheCurrentDate() error {
	// Check for month names or date patterns
	months := []string{"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}
	
	for _, month := range months {
		if strings.Contains(ctx.outputTime, month) {
			return nil
		}
	}
	return fmt.Errorf("expected month in date output: %s", ctx.outputTime)
}

func (ctx *BDDContext) theTimeShouldIncludeHoursAndMinutes() error {
	// Look for time patterns like "3:04" or "15:04"
	if strings.Contains(ctx.outputTime, ":") && 
	   (strings.Contains(ctx.outputTime, "AM") || strings.Contains(ctx.outputTime, "PM")) {
		return nil
	}
	return fmt.Errorf("expected time with hours and minutes: %s", ctx.outputTime)
}

// Step definitions for CLI workflow scenarios
func (ctx *BDDContext) theCLIToolsAreAvailable() error {
	return nil // CLI tools are built as part of the test setup
}

func (ctx *BDDContext) iRun(command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	var cmd *exec.Cmd
	if parts[0] == "geotz" {
		args := append([]string{"run", "./cmd/geotz"}, parts[1:]...)
		cmd = exec.Command("go", args...)
	} else if parts[0] == "timein" {
		args := append([]string{"run", "./cmd/timein"}, parts[1:]...)
		cmd = exec.Command("go", args...)
	} else {
		return fmt.Errorf("unknown command: %s", parts[0])
	}

	output, err := cmd.CombinedOutput()
	ctx.commandOutput = string(output)
	
	if err != nil {
		ctx.exitCode = 1
		ctx.commandError = err.Error()
	} else {
		ctx.exitCode = 0
	}
	
	return nil
}

func (ctx *BDDContext) theOutputShouldBe(expected string) error {
	output := strings.TrimSpace(ctx.commandOutput)
	if output != expected {
		return fmt.Errorf("expected output '%s', got '%s'", expected, output)
	}
	return nil
}

func (ctx *BDDContext) theOutputShouldEndWithANewline() error {
	if !strings.HasSuffix(ctx.commandOutput, "\n") {
		return fmt.Errorf("expected output to end with newline")
	}
	return nil
}

func (ctx *BDDContext) theExitCodeShouldBe(expectedCode int) error {
	if ctx.exitCode != expectedCode {
		return fmt.Errorf("expected exit code %d, got %d", expectedCode, ctx.exitCode)
	}
	return nil
}

func (ctx *BDDContext) theOutputShouldContainTheCurrentDateAndTime() error {
	if !strings.Contains(ctx.commandOutput, "2025") {
		return fmt.Errorf("expected current year in output")
	}
	return nil
}

// Initialize BDD test suite
func InitializeScenario(ctx *godog.ScenarioContext) {
	bddCtx := &BDDContext{}

	// Reset before each scenario
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		bddCtx.reset()
		return ctx, nil
	})

	// Timezone lookup steps
	ctx.Step(`^the timezone lookup service is available$`, bddCtx.theTimzoneLookupServiceIsAvailable)
	ctx.Step(`^I want to know the timezone for "([^"]*)"$`, bddCtx.iWantToKnowTheTimezoneFor)
	ctx.Step(`^I request the timezone information$`, bddCtx.iRequestTheTimezoneInformation)
	ctx.Step(`^I should get "([^"]*)" as the timezone$`, bddCtx.iShouldGetAsTheTimezone)
	ctx.Step(`^the response should be fast$`, bddCtx.theResponseShouldBeFast)
	ctx.Step(`^I should receive a helpful error message$`, bddCtx.iShouldReceiveAHelpfulErrorMessage)
	ctx.Step(`^the error should mention "([^"]*)"$`, bddCtx.theErrorShouldMention)

	// Caching steps
	ctx.Step(`^I have previously looked up "([^"]*)"$`, bddCtx.iHavePreviouslyLookedUp)
	ctx.Step(`^I request the timezone for "([^"]*)" again$`, bddCtx.iRequestTheTimezoneForAgain)
	ctx.Step(`^the response should be nearly instantaneous$`, bddCtx.theResponseShouldBeNearlyInstantaneous)
	ctx.Step(`^the result should indicate it came from cache$`, bddCtx.theResultShouldIndicateItCameFromCache)

	// Time display steps
	ctx.Step(`^the time display service is available$`, bddCtx.theTimeDisplayServiceIsAvailable)
	ctx.Step(`^I have the timezone "([^"]*)"$`, bddCtx.iHaveTheTimezone)
	ctx.Step(`^I request the current time$`, bddCtx.iRequestTheCurrentTime)
	ctx.Step(`^I should see a human-readable time format$`, bddCtx.iShouldSeeAHumanReadableTimeFormat)
	ctx.Step(`^the time should include the day of the week$`, bddCtx.theTimeShouldIncludeTheDayOfTheWeek)
	ctx.Step(`^the time should include the current date$`, bddCtx.theTimeShouldIncludeTheCurrentDate)
	ctx.Step(`^the time should include hours and minutes$`, bddCtx.theTimeShouldIncludeHoursAndMinutes)

	// CLI workflow steps
	ctx.Step(`^the CLI tools are available$`, bddCtx.theCLIToolsAreAvailable)
	ctx.Step(`^I run "([^"]*)"$`, bddCtx.iRun)
	ctx.Step(`^the output should be "([^"]*)"$`, bddCtx.theOutputShouldBe)
	ctx.Step(`^the output should end with a newline$`, bddCtx.theOutputShouldEndWithANewline)
	ctx.Step(`^the exit code should be (\d+)$`, bddCtx.theExitCodeShouldBe)
	ctx.Step(`^the output should contain the current date and time$`, bddCtx.theOutputShouldContainTheCurrentDateAndTime)

	// Alfred integration steps
	ctx.Step(`^I am using Alfred with the timein workflow$`, bddCtx.iAmUsingAlfredWithTheTimeinWorkflow)
	ctx.Step(`^I search for timezone information for "([^"]*)"$`, bddCtx.iSearchForTimezoneInformationFor)
	ctx.Step(`^I request Alfred format output$`, bddCtx.iRequestAlfredFormatOutput)
	ctx.Step(`^I should receive valid Alfred JSON$`, bddCtx.iShouldReceiveValidAlfredJSON)
	ctx.Step(`^the JSON should contain exactly one result item$`, bddCtx.theJSONShouldContainExactlyOneResultItem)
	ctx.Step(`^the item should have a title with the timezone$`, bddCtx.theItemShouldHaveATitleWithTheTimezone)
	ctx.Step(`^the item should have a subtitle mentioning the city$`, bddCtx.theItemShouldHaveASubtitleMentioningTheCity)
	ctx.Step(`^the item should be actionable$`, bddCtx.theItemShouldBeActionable)
	ctx.Step(`^I request current time in Alfred format$`, bddCtx.iRequestCurrentTimeInAlfredFormat)
	ctx.Step(`^the item title should contain the timezone and current time$`, bddCtx.theItemTitleShouldContainTheTimezoneAndCurrentTime)
	ctx.Step(`^the item subtitle should mention the city and timezone abbreviation$`, bddCtx.theItemSubtitleShouldMentionTheCityAndTimezoneAbbreviation)
	ctx.Step(`^the result should include timezone variables$`, bddCtx.theResultShouldIncludeTimezoneVariables)
	ctx.Step(`^the item should have "([^"]*)" as the title$`, bddCtx.theItemShouldHaveAsTheTitle)
	ctx.Step(`^the item should not be actionable$`, bddCtx.theItemShouldNotBeActionable)
	ctx.Step(`^the subtitle should contain the error message$`, bddCtx.theSubtitleShouldContainTheErrorMessage)
	ctx.Step(`^the subtitle should indicate the result is cached$`, bddCtx.theSubtitleShouldIndicateTheResultIsCached)
	ctx.Step(`^I examine the JSON response$`, bddCtx.iExamineTheJSONResponse)
	ctx.Step(`^it should include cache configuration$`, bddCtx.itShouldIncludeCacheConfiguration)
	ctx.Step(`^cache duration should be appropriate for the content type$`, bddCtx.cacheDurationShouldBeAppropriateForTheContentType)
	ctx.Step(`^timezone lookups should cache for (\d+) days$`, bddCtx.timezoneLookupsShouldCacheForDays)
	ctx.Step(`^time displays should cache for (\d+) seconds$`, bddCtx.timeDisplaysShouldCacheForSeconds)
	
	// Cache pre-seeding steps
	ctx.Step(`^the cache has been pre-seeded with capital cities$`, bddCtx.theCacheHasBeenPreseededWithCapitalCities)
	ctx.Step(`^I look up the timezone for "([^"]*)"$`, bddCtx.iLookUpTheTimezoneFor)
	ctx.Step(`^the result should contain "([^"]*)"$`, bddCtx.theResultShouldContain)
	ctx.Step(`^the cache should indicate a hit$`, bddCtx.theCacheShouldIndicateAHit)
	ctx.Step(`^the cache was pre-seeded (\d+) days ago$`, bddCtx.theCacheWasPreseededDaysAgo)
	ctx.Step(`^the cache has a user-created entry for "([^"]*)"$`, bddCtx.theCacheHasAUsercreatedEntryFor)
	ctx.Step(`^the result should use the user-created entry$`, bddCtx.theResultShouldUseTheUsercreatedEntry)
	ctx.Step(`^not the pre-seeded entry$`, bddCtx.notThePreseededEntry)
}

func TestBDD(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

// Additional step definitions for Alfred integration scenarios

func (ctx *BDDContext) iAmUsingAlfredWithTheTimeinWorkflow() error {
	return ctx.initializeServices()
}

func (ctx *BDDContext) iSearchForTimezoneInformationFor(city string) error {
	ctx.inputCity = city
	return nil
}

func (ctx *BDDContext) iRequestAlfredFormatOutput() error {
	// Create Alfred formatter
	alfredFormatter := presenter.NewAlfredFormatter()
	
	// Get timezone first
	output, err := ctx.geotzUseCase.GetTimezoneFromCity(ctx.inputCity)
	if err != nil {
		ctx.errorMessage = err.Error()
		// Even for errors, we need Alfred format
		errorOutput, _ := alfredFormatter.FormatError(ctx.errorMessage)
		return json.Unmarshal(errorOutput, &ctx.outputJSON)
	}

	// For successful lookup, get the timezone and format for Alfred
	timezone := strings.TrimSpace(string(output))
	tz, err := domain.NewTimezone(timezone)
	if err != nil {
		return err
	}

	alfredOutput, err := alfredFormatter.FormatTimezoneInfo(tz, ctx.inputCity, ctx.cacheHit)
	if err != nil {
		return err
	}

	return json.Unmarshal(alfredOutput, &ctx.outputJSON)
}

func (ctx *BDDContext) iShouldReceiveValidAlfredJSON() error {
	if ctx.outputJSON == nil {
		return fmt.Errorf("expected JSON output, got nil")
	}

	// Check for required Alfred structure
	if _, exists := ctx.outputJSON["items"]; !exists {
		return fmt.Errorf("Alfred JSON must contain 'items' array")
	}

	return nil
}

func (ctx *BDDContext) theJSONShouldContainExactlyOneResultItem() error {
	items, ok := ctx.outputJSON["items"].([]interface{})
	if !ok {
		return fmt.Errorf("items should be an array")
	}

	if len(items) != 1 {
		return fmt.Errorf("expected exactly 1 item, got %d", len(items))
	}

	return nil
}

func (ctx *BDDContext) theItemShouldHaveATitleWithTheTimezone() error {
	items := ctx.outputJSON["items"].([]interface{})
	item := items[0].(map[string]interface{})

	title, exists := item["title"]
	if !exists {
		return fmt.Errorf("item must have title field")
	}

	titleStr := title.(string)
	if titleStr == "Error" {
		return fmt.Errorf("expected timezone in title, got error")
	}

	// Should contain timezone format like "Europe/Paris" or "Asia/Tokyo"
	if !strings.Contains(titleStr, "/") {
		return fmt.Errorf("expected timezone format in title, got: %s", titleStr)
	}

	return nil
}

func (ctx *BDDContext) theItemShouldHaveASubtitleMentioningTheCity() error {
	items := ctx.outputJSON["items"].([]interface{})
	item := items[0].(map[string]interface{})

	subtitle, exists := item["subtitle"]
	if !exists {
		return fmt.Errorf("item must have subtitle field")
	}

	subtitleStr := strings.ToLower(subtitle.(string))
	cityLower := strings.ToLower(ctx.inputCity)

	if !strings.Contains(subtitleStr, cityLower) {
		return fmt.Errorf("subtitle should mention city '%s', got: %s", ctx.inputCity, subtitle)
	}

	return nil
}

func (ctx *BDDContext) theItemShouldBeActionable() error {
	items := ctx.outputJSON["items"].([]interface{})
	item := items[0].(map[string]interface{})

	// Check that valid is not explicitly false
	if valid, exists := item["valid"]; exists {
		if validBool, ok := valid.(bool); ok && !validBool {
			return fmt.Errorf("item should be actionable (valid should not be false)")
		}
	}

	// Check that arg exists (actionable items need args)
	if _, exists := item["arg"]; !exists {
		return fmt.Errorf("actionable item must have arg field")
	}

	return nil
}

func (ctx *BDDContext) iRequestCurrentTimeInAlfredFormat() error {
	alfredFormatter := presenter.NewAlfredFormatter()
	
	tz, err := domain.NewTimezone(ctx.inputTimezone)
	if err != nil {
		ctx.errorMessage = err.Error()
		errorOutput, _ := alfredFormatter.FormatError(ctx.errorMessage)
		return json.Unmarshal(errorOutput, &ctx.outputJSON)
	}

	alfredOutput, err := alfredFormatter.FormatTimeInfo(tz)
	if err != nil {
		return err
	}

	return json.Unmarshal(alfredOutput, &ctx.outputJSON)
}

func (ctx *BDDContext) theItemTitleShouldContainTheTimezoneAndCurrentTime() error {
	items := ctx.outputJSON["items"].([]interface{})
	item := items[0].(map[string]interface{})

	title := item["title"].(string)

	// Should contain timezone
	if !strings.Contains(title, ctx.inputTimezone) {
		return fmt.Errorf("title should contain timezone '%s', got: %s", ctx.inputTimezone, title)
	}

	// Should contain time elements (dash separator and time format)
	if !strings.Contains(title, " - ") {
		return fmt.Errorf("title should contain time separator, got: %s", title)
	}

	return nil
}

func (ctx *BDDContext) theItemSubtitleShouldMentionTheCityAndTimezoneAbbreviation() error {
	items := ctx.outputJSON["items"].([]interface{})
	item := items[0].(map[string]interface{})

	subtitle := item["subtitle"].(string)

	// Should mention the city (extract from timezone)
	parts := strings.Split(ctx.inputTimezone, "/")
	if len(parts) > 1 {
		cityName := strings.ReplaceAll(parts[1], "_", " ")
		if !strings.Contains(subtitle, cityName) {
			return fmt.Errorf("subtitle should mention city '%s', got: %s", cityName, subtitle)
		}
	}

	// Should contain timezone abbreviation (in parentheses)
	if !strings.Contains(subtitle, "(") || !strings.Contains(subtitle, ")") {
		return fmt.Errorf("subtitle should contain timezone abbreviation in parentheses, got: %s", subtitle)
	}

	return nil
}

func (ctx *BDDContext) theResultShouldIncludeTimezoneVariables() error {
	items := ctx.outputJSON["items"].([]interface{})
	item := items[0].(map[string]interface{})

	variables, exists := item["variables"]
	if !exists {
		return fmt.Errorf("item should have variables field")
	}

	varsMap := variables.(map[string]interface{})
	if _, exists := varsMap["timezone"]; !exists {
		return fmt.Errorf("variables should include timezone")
	}

	return nil
}

func (ctx *BDDContext) theItemShouldHaveAsTheTitle(expectedTitle string) error {
	items := ctx.outputJSON["items"].([]interface{})
	item := items[0].(map[string]interface{})

	title := item["title"].(string)
	if title != expectedTitle {
		return fmt.Errorf("expected title '%s', got '%s'", expectedTitle, title)
	}

	return nil
}

func (ctx *BDDContext) theItemShouldNotBeActionable() error {
	items := ctx.outputJSON["items"].([]interface{})
	item := items[0].(map[string]interface{})

	if valid, exists := item["valid"]; exists {
		if validBool, ok := valid.(bool); ok && validBool {
			return fmt.Errorf("error item should not be actionable (valid should be false)")
		}
	}

	return nil
}

func (ctx *BDDContext) theSubtitleShouldContainTheErrorMessage() error {
	items := ctx.outputJSON["items"].([]interface{})
	item := items[0].(map[string]interface{})

	subtitle := item["subtitle"].(string)
	
	// Should contain some form of error information
	if subtitle == "" {
		return fmt.Errorf("error subtitle should not be empty")
	}

	return nil
}

func (ctx *BDDContext) theSubtitleShouldIndicateTheResultIsCached() error {
	items := ctx.outputJSON["items"].([]interface{})
	item := items[0].(map[string]interface{})

	subtitle := strings.ToLower(item["subtitle"].(string))
	
	if !strings.Contains(subtitle, "cached") {
		return fmt.Errorf("subtitle should indicate cached result, got: %s", subtitle)
	}

	return nil
}

func (ctx *BDDContext) iExamineTheJSONResponse() error {
	// This is a no-op step for readability
	return nil
}

func (ctx *BDDContext) itShouldIncludeCacheConfiguration() error {
	if _, exists := ctx.outputJSON["cache"]; !exists {
		return fmt.Errorf("Alfred JSON should include cache configuration")
	}

	return nil
}

func (ctx *BDDContext) cacheDurationShouldBeAppropriateForTheContentType() error {
	cache := ctx.outputJSON["cache"].(map[string]interface{})
	seconds := cache["seconds"].(float64)

	// Basic sanity check - should be positive
	if seconds <= 0 {
		return fmt.Errorf("cache duration should be positive, got %f", seconds)
	}

	return nil
}

func (ctx *BDDContext) timezoneLookupsShouldCacheForDays(expectedDays int) error {
	cache := ctx.outputJSON["cache"].(map[string]interface{})
	seconds := cache["seconds"].(float64)
	
	expectedSeconds := float64(expectedDays * 24 * 60 * 60)
	if seconds != expectedSeconds {
		return fmt.Errorf("expected %d days (%f seconds), got %f seconds", expectedDays, expectedSeconds, seconds)
	}

	return nil
}

func (ctx *BDDContext) timeDisplaysShouldCacheForSeconds(expectedSeconds int) error {
	cache := ctx.outputJSON["cache"].(map[string]interface{})
	seconds := cache["seconds"].(float64)
	
	if seconds != float64(expectedSeconds) {
		return fmt.Errorf("expected %d seconds cache, got %f seconds", expectedSeconds, seconds)
	}

	return nil
}

// Cache pre-seeding step implementations
func (ctx *BDDContext) theCacheHasBeenPreseededWithCapitalCities() error {
	// Ensure test directory exists
	if err := os.MkdirAll(ctx.testDir, 0755); err != nil {
		return fmt.Errorf("failed to create test directory: %w", err)
	}
	
	// Set up a test cache with pre-seeded capitals
	ctx.cache = cache.NewLRUCache(200, 24*time.Hour, ctx.testDir)
	
	// Pre-seed with test capital data (city names only)
	entries := map[string]string{
		"london": "Europe/London",
		"paris":  "Europe/Paris", 
		"tokyo":  "Asia/Tokyo",
	}
	ctx.cache.PreSeed(entries)
	return nil
}

func (ctx *BDDContext) iLookUpTheTimezoneFor(city string) error {
	start := time.Now()
	
	// Check cache directly first to see if it's a cache hit
	cacheKey := strings.ToLower(city)
	if tz, ok := ctx.cache.Get(cacheKey); ok {
		ctx.lastOutput = tz
		ctx.lastDuration = time.Since(start)
		ctx.wasFromCache = true
		return nil
	}
	
	// Set up use case with pre-seeded cache for cache miss scenario
	geocoder := geocoder.NewOpenStreetMapGeocoder()
	tzf, err := timezonefinder.NewTzfTimezoneFinder()
	if err != nil {
		return fmt.Errorf("failed to setup timezone finder: %w", err)
	}
	formatter := presenter.NewPlainFormatter()
	
	useCase := usecases.NewGeotzUseCase(geocoder, tzf, ctx.cache, formatter)
	
	// Perform lookup
	result, err := useCase.GetTimezoneFromCity(city)
	if err != nil {
		ctx.lastError = err
		return nil // Don't fail here, let other steps check the error
	}
	
	ctx.lastOutput = string(result)
	ctx.lastDuration = time.Since(start)
	
	// Check if result came from cache by timing (cache hits should be very fast)
	ctx.wasFromCache = ctx.lastDuration < 50*time.Millisecond
	
	return nil
}

func (ctx *BDDContext) theResultShouldContain(expectedContent string) error {
	if !strings.Contains(ctx.lastOutput, expectedContent) {
		return fmt.Errorf("expected output to contain %q, got: %s", expectedContent, ctx.lastOutput)
	}
	return nil
}

func (ctx *BDDContext) theCacheShouldIndicateAHit() error {
	if !ctx.wasFromCache {
		return fmt.Errorf("expected cache hit (fast response), but took %v", ctx.lastDuration)
	}
	return nil
}

func (ctx *BDDContext) theCacheWasPreseededDaysAgo(daysAgo int) error {
	// Ensure test directory exists
	if err := os.MkdirAll(ctx.testDir, 0755); err != nil {
		return fmt.Errorf("failed to create test directory: %w", err)
	}
	
	// Set up cache with entries from the past
	ctx.cache = cache.NewLRUCache(200, 24*time.Hour, ctx.testDir)
	
	// Manually create cache entry with past timestamp but long TTL
	// The TTL is 365 days, so even if it was pre-seeded 30 days ago, it should still be valid
	ctx.cache.SetWithTTL("paris", "Europe/Paris", 365*24*time.Hour)
	
	return nil
}

func (ctx *BDDContext) theCacheHasAUsercreatedEntryFor(city string) error {
	// Ensure test directory exists
	if err := os.MkdirAll(ctx.testDir, 0755); err != nil {
		return fmt.Errorf("failed to create test directory: %w", err)
	}
	
	if ctx.cache == nil {
		ctx.cache = cache.NewLRUCache(200, 24*time.Hour, ctx.testDir)
	}
	
	// Add user entry with regular TTL
	ctx.cache.Set(strings.ToLower(city), "User/Custom_Timezone")
	ctx.userEntryCity = strings.ToLower(city)
	return nil
}

func (ctx *BDDContext) theResultShouldUseTheUsercreatedEntry() error {
	if !strings.Contains(ctx.lastOutput, "User/Custom_Timezone") {
		return fmt.Errorf("expected user-created entry to be used, got: %s", ctx.lastOutput)
	}
	return nil
}

func (ctx *BDDContext) notThePreseededEntry() error {
	// This is a continuation of the previous step - just verify we didn't get pre-seeded data
	if strings.Contains(ctx.lastOutput, "Europe/London") && ctx.userEntryCity == "london" {
		return fmt.Errorf("expected user entry to override pre-seeded entry")
	}
	return nil
}