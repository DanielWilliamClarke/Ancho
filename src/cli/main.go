package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/seehuhn/mt19937"

	"dwc.com/ancho/deveui"
)

func main() {
	// Setup
	maxInFlightRequests := 10
	totalDevEUIs := 100

	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())

	generator := deveui.NewDevEUIGenerator(rng)
	registrationAPI := deveui.NewRegistrationClientAPI()
	routines := deveui.NewRegistrationRoutines(registrationAPI, generator)
	batcher := deveui.NewRegistrationBatcher(routines, maxInFlightRequests)

	// Register
	registeredDevEUIs, errors := batcher.RegisterInParallel(totalDevEUIs)

	// Output
	fmt.Printf("%d DevEUIs failed ---------------- \n", len(errors))
	for _, err := range errors {
		log.Println(err)
	}
	fmt.Println("---------------------------------------------------")

	fmt.Printf("%d DevEUIs registered successfully ---------------- \n", len(registeredDevEUIs))
	for _, devEUI := range registeredDevEUIs {
		fmt.Println(devEUI)
	}
	fmt.Println("---------------------------------------------------")
}
