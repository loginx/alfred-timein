package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/loginx/alfred-timein/internal/adapters/presenter"
	"github.com/loginx/alfred-timein/internal/usecases"
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

	// Create formatter based on output format
	var formatter usecases.OutputFormatter
	if *format == "alfred" {
		formatter = presenter.NewAlfredFormatter()
	} else {
		formatter = presenter.NewPlainFormatter()
	}

	// Create use case and execute
	timeinUC := usecases.NewTimeinUseCase(formatter)
	output, err := timeinUC.GetTimezoneInfo(tz)
	if err != nil {
		outputError(err.Error(), *format)
		os.Exit(1)
	}

	os.Stdout.Write(output)
}

func outputError(msg, format string) {
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
