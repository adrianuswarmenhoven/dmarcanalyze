package main

import (
	"log/slog"
	"os"
	"time"
)

func main() {
	var programLevel = new(slog.LevelVar) // Info by default

	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))

	programLevel.Set(slog.LevelDebug)

	reps, err := getReportsViaIMAP4(Configuration.IMAP.Address+":"+Configuration.IMAP.Port, Configuration.IMAP.Username, Configuration.IMAP.Password, time.Now().AddDate(-100, 0, -1), time.Now())
	if err != nil {
		slog.Error("error getting reports", "error", err)
		os.Exit(1)
	}

	err = storeReports(reps)
	if err != nil {
		slog.Error("error storing reports", "error", err)
		os.Exit(1)
	}

}