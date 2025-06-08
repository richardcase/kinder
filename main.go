package main

import (
	"log"
	"runtime"

	"github.com/richardcase/kinder/internal"
)

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	rootCmd := internal.NewRootCommand()

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
