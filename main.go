package main

import "log"

func main() {
	fs, err := newFileStore("/home/mcetin/Dev/home-go/x/storage")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fs)

	// go func() {
	// 	log.Fatal(http.ListenAndServe("localhost:8080", globalRouter))
	// }()

	// ep, err := newEndPoint("https://ounass.ae/asdasda/asdads", "GET", true)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// ch, err := newChecker(ep, 2*time.Second, 2*time.Second)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var fn = func(ep *EndPoint) error {
	// 	log.Printf("checking \n%s", ep)
	// 	log.Println(call(ep, time.Second*2))
	// 	return nil
	// }

	// ch.start(fn)

	// time.Sleep(10 * time.Second)

	// ch.stop()

	// time.Sleep(10 * time.Second)

	// ch.start(fn)

	// time.Sleep(10 * time.Second)
}
