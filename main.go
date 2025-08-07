package main

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	loc, err := time.LoadLocation("Local")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load local time zone")
	}

	now := time.Now().In(loc)
	currentYear := now.Year()
	currentMonth := now.Month()

	ctx := context.Background()
	client := getClient(ctx)

	user, err := client.CurrentUser(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get current user")
	}

	log.Info().Msgf("Logged in as: %s (%s)", user.DisplayName, user.ID)

	likedTracks, err := fetchAllLikes(ctx, client)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to fetch liked tracks")
	}

	for _, track := range likedTracks {
		log.Debug().Msgf("Track: %s - Album: %s - Added At: %s", track.Name, track.Album.Name, track.AddedAt)
	}

	per := monthPeriod(loc, currentYear, currentMonth)
	err = runPeriod(ctx, client, per, true)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to run period for %s", per.Name)
	}
	log.Info().Msgf("Successfully processed period: %s", per.Name)
}
