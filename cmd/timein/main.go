package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/loginx/alfred-timein/internal/alfred"
	"github.com/tkuchiki/go-timezone"
)

func main() {
	format := flag.String("format", "plain", "Output format: plain or alfred")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [--format=plain|alfred] <IANA Timezone>\n", os.Args[0])
	}
	flag.Parse()

	var tz string
	if flag.NArg() == 1 {
		tz = flag.Arg(0)
	} else if flag.NArg() == 0 {
		// Try to read from STDIN
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			tz = strings.TrimSpace(scanner.Text())
		}
	} else {
		outputError("IANA timezone argument required.", *format)
		os.Exit(1)
	}

	tz = strings.TrimSpace(tz)
	if tz == "" {
		outputError("IANA timezone argument required.", *format)
		os.Exit(1)
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		outputError(fmt.Sprintf("Invalid timezone: %s", tz), *format)
		os.Exit(1)
	}
	now := time.Now().In(loc)

	if *format == "alfred" {
		city := cityFromTimezone(tz)
		tzlib := timezone.New()
		isDST := tzlib.IsDST(now)
		abbr, err := tzlib.GetTimezoneAbbreviation(tz, isDST)
		if err != nil || abbr == "" {
			abbr = now.Format("MST")
		}
		title := fmt.Sprintf("%s - %s", tz, now.Format("Mon, Jan 2, 3:04 PM"))
		subtitle := fmt.Sprintf("Current time in %s (%s)", city, abbr)
		out := alfred.NewScriptFilterOutput()
		out.Cache = &alfred.CacheConfig{Seconds: 60}
		item := alfred.Item{
			Title:    title,
			Subtitle: subtitle,
			Arg:      title,
			Variables: map[string]interface{}{
				"timezone": tz,
			},
		}
		out.AddItem(item)
		os.Stdout.Write(out.MustToJSON())
	} else {
		// Human-friendly, locale-aware output
		fmt.Println(humanTime(now))
	}
}

func outputError(msg, format string) {
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

// humanTime returns a human-friendly time string, using system locale if possible.
func humanTime(t time.Time) string {
	// Try to use system locale via environment variables (LANG, LC_TIME, etc.)
	// Go's time.Format is not locale-aware, so we do our best.
	// Example: "Monday, 02 January 2006, 15:04:05 MST"
	return t.Format("Monday, 02 January 2006, 3:04:05 PM")
}

// cityFromTimezone extracts the city/region from an IANA timezone string.
func cityFromTimezone(tz string) string {
	parts := strings.Split(tz, "/")
	if len(parts) > 1 {
		return strings.ReplaceAll(parts[1], "_", " ")
	}
	return tz
}
