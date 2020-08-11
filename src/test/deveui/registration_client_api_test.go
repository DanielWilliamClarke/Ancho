package test

import (
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/seehuhn/mt19937"

	"dwc.com/ancho/deveui"
)

type registrationTestUtils struct {
	api       deveui.IRegistrationClientAPI
	generator deveui.IDevEUIGenerator
}

func SetupRegistrationTests(url string) registrationTestUtils {
	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())
	client := &http.Client{}
	return registrationTestUtils{
		api:       deveui.LoRaWANClientAPI{client, url},
		generator: deveui.NewDevEUIGenerator(rng),
	}
}

func SetupRegistrationTestsWithNew() registrationTestUtils {
	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())
	return registrationTestUtils{
		api:       deveui.NewRegistrationClientAPI(),
		generator: deveui.NewDevEUIGenerator(rng),
	}
}

func Test_Integration_HitRegistrationEndpoint(t *testing.T) {

	// the aim of this test is to ensure that the client api can hit the endpoint successfully

	utils := SetupRegistrationTestsWithNew()

	devEUI, err := utils.generator.GenerateDevEUI(16)
	if err != nil {
		t.Errorf("Cannot generate test devEUI: %v", err)
	}

	err = utils.api.Register(devEUI)
	if err != nil {
		t.Errorf("Test registration failed: %v", err)
	}
}

func Test_Integration_Are422ErrorsRandom(t *testing.T) {

	// the aim of this test was to confirm my suspicions that the API was not reliable

	utils := SetupRegistrationTestsWithNew()
	devEUI, _ := utils.generator.GenerateDevEUI(16)

	iterations := 20
	for i := 0; i < iterations; i++ {
		err := utils.api.Register(devEUI)
		if err != nil {
			t.Logf("FAIL: DevEUI %s already registered", devEUI)
		} else {
			t.Logf("SUCCESS: %s", devEUI)
		}
	}
}

func Test_RegistrationRequestCreationFailure(t *testing.T) {

	// the aim of this test is to show that the request creation can fail and is handled

	apiHost := ":::::::::"
	utils := SetupRegistrationTests(apiHost)
	devEUI, _ := utils.generator.GenerateDevEUI(16)

	err := utils.api.Register(devEUI)
	if err == nil {
		t.Error("expected error")
	}
}

func Test_RegistrationResponseStatusFailure(t *testing.T) {

	// the aim of this test is to show that the response can return and error and is handled

	apiHost := "https://some-website.com/api-take-does-not-exist"
	utils := SetupRegistrationTests(apiHost)
	devEUI, _ := utils.generator.GenerateDevEUI(16)

	err := utils.api.Register(devEUI)
	if err == nil {
		t.Error("expected error")
	}
}
