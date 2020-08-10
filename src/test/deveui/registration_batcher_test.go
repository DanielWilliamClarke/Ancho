package test

import (
	"log"
	"testing"

	"github.com/golang/mock/gomock"

	"dwc.com/ancho/deveui"
	mock_deveui "dwc.com/ancho/test/mocks"
)

func Test_CanRunRegistrationInParallel(t *testing.T) {

	// Set up mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApi := mock_deveui.NewMockIRegistrationClientAPI(ctrl)

	totalDevEUIs := 100

	mockApi.EXPECT().
		Register(gomock.AssignableToTypeOf("test-string")).
		Times(totalDevEUIs).
		Return(nil)

	//Set up batcher with mock
	maxRequests := 10
	batcher := deveui.NewRegistrationBatcher(mockApi, maxRequests)

	// Generate some uids
	generator := deveui.NewDevEUIGenerator()
	devEUIs, err := generator.GeneratDevEUIs(totalDevEUIs, 16)
	if err != nil {
		log.Fatalf("cannot generate devEUI: %v", err)
	}

	// Run
	data, errors := batcher.RegisterInParallel(devEUIs)

	// Check
	if len(errors) > 0 {
		t.Error("Did not expect an error")
	}
	if len(data) != totalDevEUIs {
		t.Error("Batcher unable to register all DevEUIs")
	}
}
