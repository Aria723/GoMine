package main

import (
	"fmt"
	"gomine"
	time2 "time"
	"runtime"
	"os"
	"path/filepath"
	"strconv"
	"flag"
)

var currentTick = 0

var stopInstantly = false

func main() {
	var startTime = time2.Now()
	if !checkRequirements() {
		return
	}
	parseFlags()
	var serverPath = scanServerPath()

	var server, err = gomine.NewServer(serverPath)
	if err != nil {
		server.GetLogger().Critical("Another instance of the server is already running.")
		return
	}

	server.Start()
	var startupTime = time2.Now().Sub(startTime)
	server.GetLogger().Info("Server startup done! Took: " + startupTime.String())

	var tickDrop = 20

	if stopInstantly {
		server.Shutdown()
	}

	for {
		var tickDuration = int(1.0 / float32(server.GetTickRate()) * 1000) * int(time2.Millisecond)
		var nextTime = time2.Now().Add(time2.Duration(tickDuration))

		server.Tick(currentTick)

		var diff = nextTime.Sub(time2.Now()).Nanoseconds()

		if diff > 0 {
			tickDrop--

			if tickDrop < 0 && server.GetTickRate() != 20 && diff > 5 * int64(time2.Millisecond) {
				server.SetTickRate(server.GetTickRate() + 1)

				server.GetLogger().Debug("Elevating tick rate to: " + strconv.Itoa(server.GetTickRate()))
			}

			time2.Sleep(time2.Duration(diff))
		} else {
			tickDrop++

			if tickDrop > 40 {
				server.SetTickRate(server.GetTickRate() - 1)
				server.GetLogger().Debug("Lowering tick rate to: " + strconv.Itoa(server.GetTickRate()))
			}
		}

		if !server.IsRunning() {
			server.Shutdown()
			break
		}

		currentTick++
	}

	server.GetLogger().ProcessQueue() // Process the queue one last time synchronously to make sure everything gets written.
}

func scanServerPath() string {
	var executable, err = os.Executable()
	if err != nil {
		panic(err)
	}
	var serverPath = filepath.Dir(executable) + "/"

	return serverPath
}

/**
 * Checks if the Go installation meets the requirements of GoMine.
 */
func checkRequirements() bool {
	var version = runtime.Version()
	if version != "go1.9.2" {
		fmt.Println("Please install the GoLang 1.9.2 release.")
		return false
	}

	return true
}

/**
 * Parses all command line flags.
 */
func parseFlags() {
	var instantStop = flag.Bool("stop-immediately", false, "instant stop")

	flag.Parse()

	stopInstantly = *instantStop
}