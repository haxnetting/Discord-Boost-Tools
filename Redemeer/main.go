package main

import (
	"Redeemer/Core/Client"
	"Redeemer/Core/Helpers"
	b "Redeemer/Core/Keyauth"
	"time"
)

var configuration, _ = Helpers.LoadSettings()
var clear map[string]func()

// * KeyAuth Application Details *//
var name = ""
var ownerid = ""
var version = "1.0"

func main() {

	Helpers.UpdateTitle("Initializing...")
	go func() {
		for {

			b.Api(name, ownerid, version) // Important to set up the API Details
			b.Init()
			b.License(configuration.License)
			time.Sleep(time.Minute * 30)
		}
	}()

	Helpers.ClearScreen()
	Client.Redeemer()
}
