package deveui

import (
	"github.com/sonyarouje/simdb/db"
)

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

func (r ResponseCache) Load(key string) (*IdempotentPayload, error) {
	idempotentPayload := &IdempotentPayload{}
	err := r.Driver.Open(IdempotentPayload{}).Where("key", "=", key).First().AsEntity(idempotentPayload)
	if err != nil {
		return nil, err
	}
	return idempotentPayload, nil
}

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
