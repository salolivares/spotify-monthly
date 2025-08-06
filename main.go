package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	client := getClient(ctx)

	user, err := client.CurrentUser(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Logged in as:", user.DisplayName)
	likedTracks, err := fetchAllLikes(ctx, client)
	if err != nil {
		fmt.Println("error:", err)
	}

	for _, track := range likedTracks {
		fmt.Println(track.Name, " - ", track.Album.Name, " (", track.AddedAt, ")")
	}
}
