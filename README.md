# Spotify Monthly

A minimal Go CLI that snapshots your Spotify liked songs into monthly and seasonal playlists. Runs locally and uses the Spotify Web API.

## Setup

Create a Spotify developer app at <https://developer.spotify.com/dashboard>. Ensure redirect URI is set to `http://127.0.0.1:8080/callback`.

## Usage

```plaintext
SPOTIFY_ID=<id> SPOTIFY_SECRET=<secret> go run . monthly
```

## Todo

- Switch to a proper CLI library
- Cache calls to spotify API. Fetching the liked songs could be heavy.
