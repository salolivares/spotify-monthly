package main

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

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
}
