package main

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zmb3/spotify/v2"
)

// Fetch all liked tracks, this might take a while depending on the number of liked tracks.
// TODO(sal): only fetch the specific range?
func fetchAllLikes(ctx context.Context, c *spotify.Client) ([]spotify.SavedTrack, error) {
	var out []spotify.SavedTrack
	var currentPage *spotify.SavedTrackPage
	var err error

	log.Info().Msg("Fetching liked tracks...")

	// 50 is the maximum number of tracks that can be fetched in one request
	// Ordered by AddedAt, descending
	currentPage, err = c.CurrentUsersTracks(ctx, spotify.Limit(50))
	if err != nil {
		return nil, err
	}

	for page := 1; ; page++ {
		out = append(out, currentPage.Tracks...)
		err = c.NextPage(ctx, currentPage)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	log.Info().Msgf("Fetched %d liked tracks.", len(out))

	return out, nil
}
