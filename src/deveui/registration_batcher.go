package deveui

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type IRegistrationBatcher interface {
	RegisterInParallel(requiredDevEUIS int) ([]string, []error)
}

func NewRegistrationBatcher(routines IRegistrationRoutines, maxRequests int) *RegistrationBatcher {
	return &RegistrationBatcher{
		waitGroup:           &sync.WaitGroup{},
		maxInFlightRequests: maxRequests,
		routines:            routines,
	}
}

type RegistrationBatcher struct {
	waitGroup           *sync.WaitGroup
	maxInFlightRequests int
	routines            IRegistrationRoutines
}

// RegisterInParallel registers DevEUIS with the registration API in parallel, exactly requiredDevEUIS devEUIs will be registered
func (r RegistrationBatcher) RegisterInParallel(requiredDevEUIS int) ([]string, []error) {

	// Setup syscall channels
	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)
	shutdownChannel := make(chan struct{})

	// Setup data channels
	registeredChannel := make(chan string, requiredDevEUIS)
	errorChannel := make(chan error, requiredDevEUIS)

	// Run Regestration goroutines
	for index := 0; index < r.maxInFlightRequests; index++ {
		r.waitGroup.Add(1)
		go r.routines.RunBatch(ParallelRegistrationConfig{
			RegisteredChannel: registeredChannel,
			RequiredDevEUIS:   requiredDevEUIS / r.maxInFlightRequests,
			ShutdownChannel:   shutdownChannel,
			ErrorChannel:      errorChannel,
			WaitGroup:         r.waitGroup,
			GoroutineIndex:    index,
		})
	}

	// Observe total registered DevEUIs
	go r.routines.Observe(registeredChannel, requiredDevEUIS, shutdownChannel)
	// Allow SIGINT to trigger shutdown, but dont block waitgroup from joining
	go r.routines.GracefulShutdown(sigChannel, shutdownChannel)
	// cleanup observation and shutdown goroutines
	defer r.routines.CleanUp(sigChannel)

	// Wait for all waitgroups to gracefully complete
	r.waitGroup.Wait()
	fmt.Println("Registration Complete!")

	return r.processDataChannels(registeredChannel, errorChannel)
}

// processDataChannels extracts values from the data channels and returns slices
func (r RegistrationBatcher) processDataChannels(registeredChannel chan string, errorChannel chan error) ([]string, []error) {
	close(errorChannel)
	close(registeredChannel)

	registrationErrors := make([]error, 0)
	for err := range errorChannel {
		registrationErrors = append(registrationErrors, err)
	}

	registeredDevEUIs := make([]string, 0)
	for devEUI := range registeredChannel {
		registeredDevEUIs = append(registeredDevEUIs, devEUI)
	}

	return registeredDevEUIs, registrationErrors
}
