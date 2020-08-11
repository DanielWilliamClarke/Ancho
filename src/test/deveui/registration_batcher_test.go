package test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"dwc.com/ancho/deveui"
	mock_deveui "dwc.com/ancho/mocks"
)

func Test_CanRunRegistrationInParallel(t *testing.T) {

	// the aim of this test is to show that we can run n number of goroutines and still generate the expected number of devEUIs

	// Set up mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApi := mock_deveui.NewMockIRegistrationClientAPI(ctrl)

	totalDevEUIs := 1000

	mockApi.EXPECT().
		Register(gomock.AssignableToTypeOf("test-string")).
		Times(totalDevEUIs).
		Return(nil)

	//Set up batcher with mock
	maxRequests := 100

	seed := time.Now().UnixNano()
	generator := deveui.NewDevEUIGenerator(seed)
	routines := deveui.NewRegistrationRoutines(mockApi, generator)
	batcher := deveui.NewRegistrationBatcher(routines, maxRequests)

	// Run
	data, errors := batcher.RegisterInParallel(totalDevEUIs)

	// Check
	if len(errors) > 0 {
		t.Error("Did not expect an error")
	}
	if len(data) != totalDevEUIs {
		t.Error("Batcher unable to register all DevEUIs")
	}
}

func Test_CanRunRegistrationInParallelWithRegistrationFailures(t *testing.T) {

	// the aim of this test is to show that we can run n number of goroutines and still generate the expected number of devEUIs with errors

	// Set up mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApi := mock_deveui.NewMockIRegistrationClientAPI(ctrl)

	totalDevEUIs := 100
	expectedErrors := 30
	errorCounter := 0

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

	//Set up batcher with mock
	maxRequests := 10

	seed := time.Now().UnixNano()
	generator := deveui.NewDevEUIGenerator(seed)
	routines := deveui.NewRegistrationRoutines(mockApi, generator)
	batcher := deveui.NewRegistrationBatcher(routines, maxRequests)

	// Run
	data, errors := batcher.RegisterInParallel(totalDevEUIs)

	// Check
	if len(errors) < expectedErrors {
		t.Error("Not all errors accounted for")
	}

	if len(data) < totalDevEUIs {
		t.Error("Not all deveuis were registered")
	}
}
