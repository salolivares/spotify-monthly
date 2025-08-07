package main

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zmb3/spotify/v2"
)

const (
	addBatch = 100 // Spotify API allows adding up to 100 tracks at a time
	pageSize = 50  // Spotify API returns up to 50 items per page
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
	currentPage, err = c.CurrentUsersTracks(ctx, spotify.Limit(pageSize))
	if err != nil {
		return nil, err
	}

	for {
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

func ensurePlaylist(ctx context.Context, c *spotify.Client, userID, name, description string, makePublic bool) (spotify.ID, error) {
	var currentPage *spotify.SimplePlaylistPage
	var err error

	currentPage, err = c.GetPlaylistsForUser(ctx, userID, spotify.Limit(pageSize))
	if err != nil {
		return "", err
	}

	for {
		for _, pl := range currentPage.Playlists {
			if strings.EqualFold(pl.Name, name) {
				return pl.ID, nil
			}
		}
		err = c.NextPage(ctx, currentPage)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			return "", err
		}
	}

	// lil bro is missing, create it
	pl, err := c.CreatePlaylistForUser(ctx, userID, name, description, makePublic, false)
	if err != nil {
		return "", err
	}

	return pl.ID, nil
}

func getExistingPlaylistTrackIDs(ctx context.Context, c *spotify.Client, plID spotify.ID) (map[spotify.ID]struct{}, error) {
	ids := make(map[spotify.ID]struct{})
	var currentPage *spotify.PlaylistItemPage
	var err error

	currentPage, err = c.GetPlaylistItems(ctx, plID, spotify.Limit(pageSize))
	if err != nil {
		return nil, err
	}

	for {
		for _, plItem := range currentPage.Items {
			ids[plItem.Track.Track.ID] = struct{}{}
		}
		err = c.NextPage(ctx, currentPage)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return ids, nil
}

func addURIsBatched(ctx context.Context, c *spotify.Client, playlistID spotify.ID, trackIDs []spotify.ID) error {
	for i := 0; i < len(trackIDs); i += addBatch {
		end := i + addBatch
		if end > len(trackIDs) {
			end = len(trackIDs)
		}
		_, err := c.AddTracksToPlaylist(ctx, playlistID, trackIDs[i:end]...)
		if err != nil {
			return err
		}
	}
	return nil
}
