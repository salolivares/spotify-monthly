package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var err error
	var per Period

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// TODO(sal): make this more robust. it fails when not passing in a subcommand.
	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	month := flag.Int("month", 0, "Month to process (1-12). Defaults to current month.")

	flag.CommandLine.Init(os.Args[0]+" "+cmd, flag.ContinueOnError)
	_ = flag.CommandLine.Parse(os.Args[2:])

	loc, err := time.LoadLocation("Local")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load local time zone")
	}

	now := time.Now().In(loc)
	y := now.Year()
	m := now.Month()
	if *month != 0 {
		m = time.Month(*month)
	}

	ctx := context.Background()
	client := getClient(ctx)

	switch cmd {
	case "monthly":
		per = monthPeriod(loc, y, m)
		err = runPeriod(ctx, client, per, true)
	default:
		fmt.Println("Usage: spotifyx monthly [--month <1-12>]")
		os.Exit(1)
	}

	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to run period for %s", per.Name)
	}

	log.Info().Msgf("Successfully processed period: %s", per.Name)
}
