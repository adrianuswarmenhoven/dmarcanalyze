package main

import (
	"log/slog"
	"os"
	"time"
)

func main() {
fetchLoop:
	for {
		startTime, err := getLastRun()
		if err != nil {
			slog.Error("error getting last run", "error", err)
			os.Exit(1)
		}
		setLastRun(time.Now().Add(-time.Second))
		endTime := time.Now().AddDate(100, 0, 1)
		timerFetch := time.Now()
		reps, err := getReportsViaIMAP4(Configuration.IMAP.Address+":"+Configuration.IMAP.Port, Configuration.IMAP.Username, Configuration.IMAP.Password, startTime, endTime)
		if err != nil {
			slog.Info("no new reports found", "error", err)
			if Configuration.Sleep == 0 {
				os.Exit(0)
			} else {
				// Sleep until the next run
				slog.Info("sleeping", "duration", Configuration.Sleep)
				time.Sleep(time.Duration(Configuration.Sleep) * time.Second)
				continue fetchLoop
			}
		}
		durationFetch := time.Since(timerFetch)

		timerStore := time.Now()
		err = storeReports(reps)
		if err != nil {
			slog.Error("error storing reports", "error", err)
			os.Exit(1)
		}
		durationStore := time.Since(timerStore)
		slog.Info("finished", "durationFetch", durationFetch, "durationStore", durationStore, "total", durationFetch+durationStore)
		if Configuration.Sleep == 0 {
			os.Exit(0)
		}
		// Sleep until the next run
		slog.Info("sleeping", "duration", Configuration.Sleep)
		time.Sleep(time.Duration(Configuration.Sleep) * time.Second)
	}
}
