package main

import (
	"log"
	"time"
)

func main() {

	ep, err := newEndPoint("http://www.google.com", "GET", true)
	if err != nil {
		log.Fatal(err)
	}

	call(ep, time.Second*2)

	// ch, err := newChecker(ep, 2*time.Second, 2*time.Second)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var fn = func(ep *endPoint) error {
	// 	log.Printf("checking \n%s", ep)
	// 	return nil
	// }

	// ch.start(fn)

	// time.Sleep(10 * time.Second)

	// ch.stop()

	// time.Sleep(10 * time.Second)

	// ch.start(fn)

	// time.Sleep(10 * time.Second)
}
