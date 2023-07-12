# Run go using nodemon to allow go to be reloaded upon changes.
nodemon --exec go run main.go --signal SIGTERM
# nodemon --exec "go run" main.go