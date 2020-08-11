package deveui

import (
	"github.com/sonyarouje/simdb/db"
)

// Setup db type
type IdempotentPayload struct {
	Key     string        `json:"key"`
	Payload DevEUIPayload `json:"payload"`
}

func (i IdempotentPayload) ID() (jsonField string, value interface{}) {
	value = i.Key
	jsonField = "key"
	return
}

type IResponseCache interface {
	Load(key string) (IdempotentPayload, error)
	Store(key string, data interface{}) error
}

func NewResponseCache() (*ResponseCache, error) {
	driver, err := db.New("idempotency_db")
	if err != nil {
		return nil, err
	}
	return &ResponseCache{
		Driver: driver,
	}, nil
}

type ResponseCache struct {
	Driver *db.Driver
}

// Load loads data from the simple db by key will return an error if no data is found
func (r ResponseCache) Load(key string) (*IdempotentPayload, error) {
	idempotentPayload := &IdempotentPayload{}
	err := r.Driver.Open(IdempotentPayload{}).Where("key", "=", key).First().AsEntity(idempotentPayload)
	if err != nil {
		return nil, err
	}
	return idempotentPayload, nil
}

// Store stores data by key, if the key already exists it will attempt to update the record instead
func (r ResponseCache) Store(key string, data DevEUIPayload) error {
	_, err := r.Load(key)
	if err == nil {
		err = r.Driver.Update(IdempotentPayload{
			Key:     key,
			Payload: data,
		})
	} else {
		err = r.Driver.Insert(IdempotentPayload{
			Key:     key,
			Payload: data,
		})
	}
	return err
}
