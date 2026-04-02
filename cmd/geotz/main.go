package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/loginx/alfred-timein/internal/adapters/cache"
	"github.com/loginx/alfred-timein/internal/adapters/geocoder"
	"github.com/loginx/alfred-timein/internal/adapters/presenter"
	"github.com/loginx/alfred-timein/internal/adapters/timezonefinder"
	"github.com/loginx/alfred-timein/internal/domain"
	"github.com/loginx/alfred-timein/internal/usecases"
)

func main() {
	format := flag.String("format", "plain", "Output format: plain or alfred")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [--format=plain|alfred] <city or landmark>\n", os.Args[0])
	}
	flag.Parse()

	if flag.NArg() < 1 {
		outputError("City or landmark argument required.", *format)
		os.Exit(1)
	}
	query := strings.Join(flag.Args(), " ")
	city := strings.TrimSpace(query)
	if city == "" {
		outputError("City or landmark argument required.", *format)
		os.Exit(1)
	}

	// Use Alfred's persistent data dir when available, fall back to CWD for standalone CLI
	cacheDir := resolveCacheDir()
	if cacheDir != "." {
		seedFromWorkflowDir(cacheDir)
	}
	cacheAdapter := cache.NewLRUCache(1000, 30*24*time.Hour, cacheDir)
	cacheKey := strings.ToLower(city)
	if tz, ok := cacheAdapter.Get(cacheKey); ok {
		// Cache hit - skip expensive validation, just format and output
		if *format == "alfred" {
			formatter := presenter.NewAlfredFormatter()
			timezone := &domain.Timezone{Name: tz}
			output, err := formatter.FormatTimezoneInfo(timezone, city, true)
			if err != nil {
				outputError(err.Error(), *format)
				os.Exit(1)
			}
			os.Stdout.Write(output)
		} else {
			// Plain format - just output timezone name directly for maximum speed
			fmt.Println(tz)
		}
		return
	}

	// Cache miss - initialize all dependencies and use full use case
	var formatter usecases.OutputFormatter
	if *format == "alfred" {
		formatter = presenter.NewAlfredFormatter()
	} else {
		formatter = presenter.NewPlainFormatter()
	}

	geocoderAdapter := geocoder.NewOpenStreetMapGeocoder()
	
	tzFinder, err := timezonefinder.NewTzfTimezoneFinder()
	if err != nil {
		outputError("Failed to initialize timezone finder.", *format)
		os.Exit(1)
	}

	// Create use case and execute
	geotzUC := usecases.NewGeotzUseCase(geocoderAdapter, tzFinder, cacheAdapter, formatter)
	output, err := geotzUC.GetTimezoneFromCity(city)
	if err != nil {
		// For plain format, write errors to stderr and exit with error code
		if *format == "plain" {
			fmt.Fprintln(os.Stderr, "Error:", err.Error())
			os.Exit(1)
		} else {
			outputError(err.Error(), *format)
			os.Exit(1)
		}
	}

	os.Stdout.Write(output)
}

// resolveCacheDir returns the directory for runtime cache storage.
// Inside Alfred: uses alfred_workflow_data (persistent across updates).
// Outside Alfred: falls back to current directory.
func resolveCacheDir() string {
	if dir := os.Getenv("alfred_workflow_data"); dir != "" {
		return dir
	}
	return "."
}

// seedFromWorkflowDir copies the shipped pre-seeded cache to the data directory
// on first run. No-op if the data dir cache already exists.
func seedFromWorkflowDir(dataDir string) {
	dest := filepath.Join(dataDir, "geotz_cache.json")
	if _, err := os.Stat(dest); err == nil {
		return
	}
	src, err := os.ReadFile(filepath.Join(".", "geotz_cache.json"))
	if err != nil {
		return
	}
	os.MkdirAll(dataDir, 0755)
	os.WriteFile(dest, src, 0644)
}

func outputError(msg string, format string) {
	var formatter usecases.OutputFormatter
	if format == "alfred" {
		formatter = presenter.NewAlfredFormatter()
	} else {
		formatter = presenter.NewPlainFormatter()
	}
	
	output, err := formatter.FormatError(msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error formatting error message:", err)
		return
	}
	
	if format == "alfred" {
		os.Stdout.Write(output)
	} else {
		fmt.Fprintln(os.Stderr, string(output))
	}
}
