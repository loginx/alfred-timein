package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/loginx/alfred-timein/internal/adapters/cache"
	"github.com/loginx/alfred-timein/internal/adapters/geocoder"
	"github.com/loginx/alfred-timein/internal/adapters/presenter"
	"github.com/loginx/alfred-timein/internal/adapters/timezonefinder"
	"github.com/loginx/alfred-timein/internal/domain"
	"github.com/loginx/alfred-timein/internal/usecases"
)

const cacheSeconds = 604800 // 7 days

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

	// Fast path: Check cache first before initializing expensive dependencies
	cacheAdapter := cache.NewDefaultCache()
	cacheKey := strings.ToLower(city)
	if tz, ok := cacheAdapter.Get(cacheKey); ok {
		// Cache hit - create minimal dependencies for formatting only
		var formatter usecases.OutputFormatter
		if *format == "alfred" {
			formatter = presenter.NewAlfredFormatter()
		} else {
			formatter = presenter.NewPlainFormatter()
		}
		
		timezone, err := domain.NewTimezone(tz)
		if err != nil {
			outputError(err.Error(), *format)
			os.Exit(1)
		}
		
		output, err := formatter.FormatTimezoneInfo(timezone, city, true)
		if err != nil {
			outputError(err.Error(), *format)
			os.Exit(1)
		}
		
		os.Stdout.Write(output)
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
