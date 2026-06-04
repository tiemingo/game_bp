package client

import (
	"sync"

	"github.com/Liphium/neoroute"
)

var clients = &sync.Map{} // map[string]*sync.Mutex

type ClientData struct {
	PlayerId string
	LobbyId  string
}

func NewClientData(sessionId string) ClientData {
	return ClientData{}
}

func AccessData(ctx *neoroute.Ctx[ClientData], modifyFunc func(*ClientData) error) error {
	mutexAny, _ := clients.LoadOrStore(ctx.Session().Id(), &sync.Mutex{})
	mutex := mutexAny.(*sync.Mutex)

	mutex.Lock()
	defer mutex.Unlock()

	data := ctx.Session().Data()
	dataToModify := &data
	err := modifyFunc(dataToModify)
	ctx.Session().SetData(*dataToModify)
	return err
}
