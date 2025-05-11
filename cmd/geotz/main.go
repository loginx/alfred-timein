package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/codingsince1985/geo-golang/openstreetmap"
	"github.com/ringsaturn/tzf"

	"github.com/loginx/alfred-timein/internal/alfred"
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

	cache := alfred.DefaultGeotzCache()
	cacheKey := strings.ToLower(city)
	if tz, ok := cache.Get(cacheKey); ok {
		outputResult(city, tz, true, *format)
		return
	}

	geocoder := openstreetmap.Geocoder()
	loc, err := geocoder.Geocode(city)
	if err != nil || loc == nil {
		outputError(fmt.Sprintf("Could not geocode: %s", city), *format)
		os.Exit(1)
	}

	finder, err := tzf.NewDefaultFinder()
	if err != nil {
		outputError("Failed to initialize timezone finder.", *format)
		os.Exit(1)
	}
	tz := finder.GetTimezoneName(float64(loc.Lng), float64(loc.Lat))
	if tz == "" {
		outputError(fmt.Sprintf("Could not resolve timezone for: %s", city), *format)
		os.Exit(1)
	}

	cache.Set(cacheKey, tz)
	outputResult(city, tz, false, *format)
}

func outputResult(city, tz string, cached bool, format string) {
	if format == "alfred" {
		out := alfred.NewScriptFilterOutput()
		out.Cache = &alfred.CacheConfig{Seconds: cacheSeconds}
		sub := city
		if cached {
			sub += " (cached)"
		}
		item := alfred.Item{
			Title:    tz,
			Subtitle: sub,
			Arg:      tz,
			Variables: map[string]interface{}{
				"city": city,
			},
		}
		out.AddItem(item)
		os.Stdout.Write(out.MustToJSON())
	} else {
		fmt.Println(tz)
	}
}

func outputError(msg string, format string) {
	if format == "alfred" {
		out := alfred.NewScriptFilterOutput()
		item := alfred.Item{
			Title:    "Error",
			Subtitle: msg,
			Valid:    boolPtr(false),
		}
		out.AddItem(item)
		os.Stdout.Write(out.MustToJSON())
	} else {
		fmt.Fprintln(os.Stderr, msg)
	}
}

func boolPtr(b bool) *bool {
	return &b
}
