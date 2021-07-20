package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func interrupt(cb func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Fprint(w, "\r")
		cb()
	}()
}

func dirExists(fp string) bool {
	stat, err := os.Stat(fp)
	return err == nil && stat.IsDir()
}

// func fileExists(fp string) bool {
// 	_, err := os.Stat(fp)
// 	return err == nil
// }
