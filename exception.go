package cache

import (
	"fmt"
	. "github.com/gopi-frame/contract/exception"
	"github.com/gopi-frame/exception"
)

type StoreNotConfiguredException struct {
	storeName string
	Throwable
}

func NewStoreNotConfiguredException(storeName string) *StoreNotConfiguredException {
	return &StoreNotConfiguredException{
		storeName: storeName,
		Throwable: exception.New(fmt.Sprintf("store [%s] not configured", storeName)),
	}
}

func (e *StoreNotConfiguredException) StoreName() string {
	return e.storeName
}
