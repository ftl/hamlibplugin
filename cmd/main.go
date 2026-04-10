package main

import (
	"log"

	sdk "github.com/SkYNewZ/streamdeck-sdk"

	"github.com/ftl/hamlibplugin/pkg/plugin"
)

var version = "development"

func main() {
	streamDeck, err := sdk.New()
	if err != nil {
		log.Fatal(err)
	}
	plugin := plugin.New(streamDeck.UUID, streamDeck)
	streamDeck.Handler(plugin)
	plugin.Start()
	streamDeck.Start()
}
