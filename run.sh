#!/bin/bash
# Run go using nodemon to allow go to be reloaded upon changes.
# Run go in realtime
nodemon --exec go run main.go --signal SIGTERM -- development
# nodemon --exec "go run" main.go