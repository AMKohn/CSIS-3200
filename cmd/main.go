package main

import (
	"csis3200/internal/app/dashboard"
	"csis3200/internal/app/ingester"
	"sync"
)

func main() {
	println("Hello world")

	// Create the waitgroup so the program keeps running until both servers exit
	var wg sync.WaitGroup
	wg.Add(2)

	dashboard.StartServer(&wg)
	ingester.StartServer(&wg)

	wg.Wait()
	println("Servers done running, exiting")
}