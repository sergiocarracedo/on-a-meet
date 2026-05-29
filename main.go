package main

import "github.com/sergiocarracedo/on-a-meet/cmd"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Execute()
}
