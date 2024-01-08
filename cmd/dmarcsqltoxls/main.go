package main

import (
	"log/slog"
	"os"
)

func main() {

	err := initDB()
	if err != nil {
		slog.Error("error initializing database", "error", err)
		os.Exit(1)
	}

	buildXLS()
}
