package test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"dwc.com/ancho/deveui"
)

func Test_CanCreateTestLocalDb(t *testing.T) {
	_, err := deveui.NewResponseCache()
	if err != nil {
		t.Errorf("Cannot create cache: %v", err)
	}
}

func Test_CanStoreIntoLocalDb(t *testing.T) {
	cache, err := deveui.NewResponseCache()
	if err != nil {
		t.Errorf("Cannot create cache: %v", err)
	}

	generator := deveui.NewDevEUIGenerator(time.Now().UnixNano())

	id, err := generator.GenerateDevEUI(16)
	if err != nil {
		log.Printf("Could not generate id: %v", err)
	}

	payload := deveui.DevEUIPayload{
		DevEUIs: []string{id, id, id},
	}

	err = cache.Store(
		fmt.Sprintf("Test-%s", id),
		payload)

	if err != nil {
		t.Errorf("Could not store idempotent payload data: %v", err)
	}
}

func Test_CanStoreAndLoadFromLocalDb(t *testing.T) {
	cache, err := deveui.NewResponseCache()
	if err != nil {
		t.Errorf("Cannot create cache: %v", err)
	}

	generator := deveui.NewDevEUIGenerator(time.Now().UnixNano())

	id, err := generator.GenerateDevEUI(16)
	if err != nil {
		log.Printf("Could not generate id: %v", err)
	}

	payload := deveui.DevEUIPayload{
		DevEUIs: []string{id, id, id},
	}

	key := fmt.Sprintf("Test-%s", id)

	err = cache.Store(key, payload)

	if err != nil {
		t.Errorf("Could not store idempotent payload data: %v", err)
	}

	_, err = cache.Load(key)
	if err != nil {
		t.Errorf("Could not retrieve idempotent payload data: %v", err)
	}
}

func Test_CanUpdateInLocalDb(t *testing.T) {
	cache, err := deveui.NewResponseCache()
	if err != nil {
		t.Errorf("Cannot create cache: %v", err)
	}

	generator := deveui.NewDevEUIGenerator(time.Now().UnixNano())

	id, err := generator.GenerateDevEUI(16)
	if err != nil {
		log.Printf("Could not generate id: %v", err)
	}
	key := fmt.Sprintf("Test-%s", id)

	// Store empty data
	err = cache.Store(key, deveui.DevEUIPayload{
		DevEUIs: []string{},
	})
	if err != nil {
		t.Errorf("Could not store idempotent payload data: %v", err)
	}
	data, err := cache.Load(key)
	if err != nil {
		t.Errorf("Could not retrieve idempotent payload data: %v", err)
	}
	if len(data.Payload.DevEUIs) > 0 {
		t.Error("Did not expect data for this key")
	}

	// Update with new data
	err = cache.Store(key, deveui.DevEUIPayload{
		DevEUIs: []string{"Now-some-data"},
	})
	if err != nil {
		t.Errorf("Could not store idempotent payload data: %v", err)
	}
	data, err = cache.Load(key)
	if err != nil {
		t.Errorf("Could not retrieve idempotent payload data: %v", err)
	}
	if len(data.Payload.DevEUIs) != 1 {
		t.Error("Did not expect data for this key")
	}
}
