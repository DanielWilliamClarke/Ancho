package test

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"dwc.com/ancho/deveui"
	mock_deveui "dwc.com/ancho/mocks"
)

func runBatchTestUtil(routines deveui.IRegistrationRoutines, waitGroup *sync.WaitGroup, shutdownChannel chan struct{}, requiredDevEUIS int) deveui.ParallelRegistrationConfig {

	registeredChannel := make(chan string, requiredDevEUIS)
	errorChannel := make(chan error, requiredDevEUIS)
	waitGroup.Add(1)

	config := deveui.ParallelRegistrationConfig{
		RegisteredChannel: registeredChannel,
		RequiredDevEUIS:   requiredDevEUIS,
		ShutdownChannel:   shutdownChannel,
		ErrorChannel:      errorChannel,
		WaitGroup:         waitGroup,
		GoroutineIndex:    1,
	}
	go routines.RunBatch(config)

	return config
}

func Test_CanRunBatchAndInterrupt(t *testing.T) {

	// The aim of this test is to show that a run batch goroutine can be quit early gracefully

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApi := mock_deveui.NewMockIRegistrationClientAPI(ctrl)
	mockGenerator := mock_deveui.NewMockIDevEUIGenerator(ctrl)

	mockGenerator.EXPECT().
		GenerateDevEUI(gomock.AssignableToTypeOf(16)).
		AnyTimes().
		Return("test-deveui", nil)

	mockApi.EXPECT().
		Register(gomock.AssignableToTypeOf("test-string")).
		AnyTimes().
		DoAndReturn(func(devEUI string) error {
			time.Sleep(2 * time.Second)
			return nil
		})

	routines := deveui.NewRegistrationRoutines(mockApi, mockGenerator)
	waitGroup := &sync.WaitGroup{}
	shutdownChannel := make(chan struct{})

	config := runBatchTestUtil(routines, waitGroup, shutdownChannel, 10)

	time.Sleep(1 * time.Second)
	close(shutdownChannel)

	waitGroup.Wait()

	close(config.ErrorChannel)
	close(config.RegisteredChannel)

	registeredDevEUIs := make([]string, 0)
	for devEUI := range config.RegisteredChannel {
		registeredDevEUIs = append(registeredDevEUIs, devEUI)
	}

	if len(registeredDevEUIs) < 0 {
		t.Error("No deveuis were registered")
	}
}

func Test_RunBatchGeneratorError(t *testing.T) {

	// The aim of this test is to show that a generator error will force the goroutine to quit

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApi := mock_deveui.NewMockIRegistrationClientAPI(ctrl)
	mockGenerator := mock_deveui.NewMockIDevEUIGenerator(ctrl)

	mockGenerator.EXPECT().
		GenerateDevEUI(gomock.AssignableToTypeOf(16)).
		Times(1).
		Return("", errors.New("fake generator error"))

	mockApi.EXPECT().
		Register(gomock.AssignableToTypeOf("test-string")).
		Times(0)

	routines := deveui.NewRegistrationRoutines(mockApi, mockGenerator)
	waitGroup := &sync.WaitGroup{}
	shutdownChannel := make(chan struct{})
	config := runBatchTestUtil(routines, waitGroup, shutdownChannel, 10)

	waitGroup.Wait()

	close(config.ErrorChannel)
	close(config.RegisteredChannel)

	registrationErrors := make([]error, 0)
	for err := range config.ErrorChannel {
		registrationErrors = append(registrationErrors, err)
	}

	if len(registrationErrors) < 0 {
		t.Error("Not all errors accounted for")
	}
}

func Test_RunBatchAPIError(t *testing.T) {

	// The aim of this test is to show that a batch goroutine will continue to run until its expected number of deveuis have been registered

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApi := mock_deveui.NewMockIRegistrationClientAPI(ctrl)
	mockGenerator := mock_deveui.NewMockIDevEUIGenerator(ctrl)

	totalDevEUIs := 10
	expectedErrors := 3
	errorCounter := 0

	mockGenerator.EXPECT().
		GenerateDevEUI(gomock.AssignableToTypeOf(16)).
		Times(totalDevEUIs+expectedErrors).
		Return("test-deveui", nil)

	mockApi.EXPECT().
		Register(gomock.AssignableToTypeOf("test-deveui")).
		Times(totalDevEUIs + expectedErrors).
		DoAndReturn(func(deveui string) error {
			if errorCounter < expectedErrors {
				errorCounter++
				return errors.New("fake api error")
			}
			return nil
		})

	routines := deveui.NewRegistrationRoutines(mockApi, mockGenerator)
	waitGroup := &sync.WaitGroup{}
	shutdownChannel := make(chan struct{})
	config := runBatchTestUtil(routines, waitGroup, shutdownChannel, totalDevEUIs)

	waitGroup.Wait()

	close(config.ErrorChannel)
	close(config.RegisteredChannel)

	registrationErrors := make([]error, 0)
	for err := range config.ErrorChannel {
		registrationErrors = append(registrationErrors, err)
	}

	registeredDevEUIs := make([]string, 0)
	for devEUI := range config.RegisteredChannel {
		registeredDevEUIs = append(registeredDevEUIs, devEUI)
	}

	if len(registrationErrors) < expectedErrors {
		t.Error("Not all errors accounted for")
	}

	if len(registeredDevEUIs) < totalDevEUIs {
		t.Error("Not all deveuis were registered")
	}
}

func Test_CanObserveRegisteredChannel(t *testing.T) {

	// The aim of this test is to show that the observation goroutine ends when all deveuis are registered

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApi := mock_deveui.NewMockIRegistrationClientAPI(ctrl)
	mockGenerator := mock_deveui.NewMockIDevEUIGenerator(ctrl)
	routines := deveui.NewRegistrationRoutines(mockApi, mockGenerator)

	// The aim of this test is to show that a run batch goroutine can be quit early gracefully
	totalDevEUIs := 3
	registeredChannel := make(chan string, totalDevEUIs)
	shutdownChannel := make(chan struct{})

	go routines.Observe(registeredChannel, totalDevEUIs, shutdownChannel)

	for len(registeredChannel) < totalDevEUIs {
		time.Sleep(1 * time.Second)
		registeredChannel <- "test"
	}
}

func Test_CanObserveRegisteredChannelAndInterrupt(t *testing.T) {

	// The aim of this test is to show that the observation goroutine can be interrupted and shutdown gracefully

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApi := mock_deveui.NewMockIRegistrationClientAPI(ctrl)
	mockGenerator := mock_deveui.NewMockIDevEUIGenerator(ctrl)
	routines := deveui.NewRegistrationRoutines(mockApi, mockGenerator)

	// The aim of this test is to show that a run batch goroutine can be quit early gracefully
	totalDevEUIs := 10
	registeredChannel := make(chan string, totalDevEUIs)
	shutdownChannel := make(chan struct{})

	go routines.Observe(registeredChannel, totalDevEUIs, shutdownChannel)

	time.Sleep(1 * time.Second)
	registeredChannel <- "test"
	time.Sleep(1 * time.Second)
	close(shutdownChannel)
}

func Test_CanHandleGracefulShutdown(t *testing.T) {

	// The aim of this test is to show that the SIGINT channel can be closed programmatically to clean up this goroutine

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApi := mock_deveui.NewMockIRegistrationClientAPI(ctrl)
	mockGenerator := mock_deveui.NewMockIDevEUIGenerator(ctrl)
	routines := deveui.NewRegistrationRoutines(mockApi, mockGenerator)

	// The aim of this test is to show that a run batch goroutine can be quit early gracefully
	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)
	shutdownChannel := make(chan struct{})

	go routines.GracefulShutdown(sigChannel, shutdownChannel)
	routines.CleanUp(sigChannel)
}
