package main

import "context"

func main() {
	ctx := context.Background()
	client := getClient(ctx)

	user, err := client.CurrentUser(ctx)
	if err != nil {
		panic(err)
	}
	println("Logged in as:", user.DisplayName)
}
