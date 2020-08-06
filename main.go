package main

import (
	"github.com/valyala/fastjson"
	"log"
	"meow.tf/streamdeck/sdk"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var context string // set after plugin registration
var lastMicOnVolume = 0

func main() {
	// Initialize handlers for events
	sdk.RegisterAction("us.hervalicio.sdmicmac", handleToggleVolume)
	sdk.AddHandler(func(e *sdk.WillAppearEvent) {
		context = e.Context
	})

	// Open and connect the SDK
	err := sdk.Open()
	if err != nil {
		log.Fatalln(err)
	}

	go pollMic()

	// Wait until the socket is closed, or SIGTERM/SIGINT is received
	sdk.Wait()
}

func execAndParseVolume(script string) int {
	cmd := exec.Command("osascript", script)

	out, err := cmd.CombinedOutput()
	if err != nil {
		sdk.Log("FAIL: " + err.Error())
		panic(err)
	}

	vol, err := strconv.Atoi(
		strings.Trim(string(out), "\n"),
	)
	if err != nil {
		sdk.Log("Failed parsing volume: " + err.Error())
		panic(err)
	}

	return vol
}

func pollMic() {
	for {
		sdk.Log("Polling mic state..." + sdk.PluginUUID)
		setMicVol(
			execAndParseVolume("mic_state.workflow"),
		)

		time.Sleep(time.Second * 5)
	}
}

func toggleMic() {
	setMicVol(
		execAndParseVolume("mic.workflow"), // TODO pass in last volume
	)
}

func setMicVol(vol int) {
	sdk.Log("Setting vol to " + strconv.Itoa(vol))

	if vol > 0 {
		lastMicOnVolume = vol
	}

	// mic is disabled
	if vol == 0 {
		sdk.SetState(context, 1)
	} else {
		sdk.SetState(context, 0)
	}
}

func handleToggleVolume(action, context string, payload *fastjson.Value, deviceId string) {
	sdk.Log(context)

	toggleMic()
}
