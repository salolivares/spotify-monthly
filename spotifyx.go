package main

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
)

// Fetch all liked tracks, this might take a while depending on the number of liked tracks.
// TODO(sal): only fetch the specific range?
func fetchAllLikes(ctx context.Context, c *spotify.Client) ([]spotify.SavedTrack, error) {
	var out []spotify.SavedTrack
	var currentPage *spotify.SavedTrackPage
	var err error

	// 50 is the maximum number of tracks that can be fetched in one request
	// Ordered by AddedAt, descending
	currentPage, err = c.CurrentUsersTracks(ctx, spotify.Limit(50))
	if err != nil {
		return nil, err
	}

	fmt.Println("Fetching liked tracks...")
	for page := 1; ; page++ {
		out = append(out, currentPage.Tracks...)
		err = c.NextPage(ctx, currentPage)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			fmt.Println("Error fetching page", page, ":", err)
		}
	}
	fmt.Println("Fetched", len(out), "liked tracks.")

	return out, nil
}
