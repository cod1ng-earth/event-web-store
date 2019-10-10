package main

import (
	"log"

	"github.com/cod1ng-earth/event-web-store/backend/pkg/simba"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {

	readLock := kingpin.Flag("lock", "lock to protect model: exclusive, parallel, wait-free").String()
	bridges := kingpin.Flag("bridge", "connections to other contexts").Strings()
	kingpin.Parse()

	if *readLock == "" {
		*readLock = "parallel"
	}
	if *readLock != "exclusive" && *readLock != "parallel" && *readLock != "wait-free" {
		log.Fatalf("readLock '%s' unknown", *readLock)
	}

	simba.UpdateCode(*readLock, *bridges)
}
