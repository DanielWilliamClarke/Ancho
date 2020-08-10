package deveui

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type IRegistrationBatcher interface {
	RegisterInParallel(devEUIs []string) ([]string, []error)
}

func NewRegistrationBatcher(api IRegistrationClientAPI, maxRequests int) *RegistrationBatcher {
	return &RegistrationBatcher{
		WaitGroup:           &sync.WaitGroup{},
		MaxInFlightRequests: maxRequests,
		API:                 api,
	}
}

type parallelRegistrationConfig struct {
	registeredChannel chan string
	devEUIs           []string
	shutdownChannel   chan struct{}
	errorChannel      chan error
	waitGroup         *sync.WaitGroup
	goroutineIndex    int
}

type RegistrationBatcher struct {
	WaitGroup           *sync.WaitGroup
	MaxInFlightRequests int
	API                 IRegistrationClientAPI
}

func (r RegistrationBatcher) RegisterInParallel(devEUIs []string) ([]string, []error) {

	// Divide devEUIs between workloads
	devEUIBatches := ChunkDevEUIs(devEUIs, len(devEUIs)/r.MaxInFlightRequests)

	// Setup syscall channels
	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)
	shutdownChannel := make(chan struct{})

	// Setup data channels
	registeredChannel := make(chan string, len(devEUIs))
	errorChannel := make(chan error, len(devEUIs))

	// Run Regestration goroutines
	for index := 0; index < r.MaxInFlightRequests; index++ {
		r.WaitGroup.Add(1)
		go r.runBatch(parallelRegistrationConfig{
			registeredChannel: registeredChannel,
			devEUIs:           devEUIBatches[index],
			shutdownChannel:   shutdownChannel,
			errorChannel:      errorChannel,
			waitGroup:         r.WaitGroup,
			goroutineIndex:    index,
		})
	}

	// Allow SIGINT to trigger shutdown, but dont block waitgroup from joining
	defer func() {
		fmt.Println("Cleaning Up")
		close(sigChannel)
	}()
	go func() {
		<-sigChannel // received SIGINT or SIGTERM
		close(shutdownChannel)
		// Wait for all waitgroups to gracefully complete
		fmt.Println("Quit signal received, gracefully shutdown registration...")
	}()

	// Wait for goroutines to complete
	r.WaitGroup.Wait()
	fmt.Println("Registration Complete!")

	return r.processDataChannels(registeredChannel, errorChannel)
}

func (r RegistrationBatcher) runBatch(config parallelRegistrationConfig) {
	fmt.Printf("Starting registration batch: %d \n", config.goroutineIndex)
	defer config.waitGroup.Done()

	for _, devEUI := range config.devEUIs {
		select {
		case <-config.shutdownChannel:
			log.Printf("Shutdown registration batch: %d \n", config.goroutineIndex)
			return
		default:
			err := r.API.Register(devEUI)
			if err != nil {
				config.errorChannel <- err
			} else {
				config.registeredChannel <- devEUI
			}
		}
	}

	fmt.Printf("Registration batch %d complete \n", config.goroutineIndex)
}

func (r RegistrationBatcher) processDataChannels(registeredChannel chan string, errorChannel chan error) ([]string, []error) {
	close(errorChannel)
	close(registeredChannel)

	registrationErrors := make([]error, 0)
	for err := range errorChannel {
		if err != nil {
			registrationErrors = append(registrationErrors, err)
		}
	}

	registeredDevEUIs := make([]string, 0)
	for devEUI := range registeredChannel {
		registeredDevEUIs = append(registeredDevEUIs, devEUI)
	}

	return registeredDevEUIs, registrationErrors
}
