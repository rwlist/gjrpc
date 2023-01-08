package transport

import (
	"errors"

	"github.com/rwlist/gjrpc/pkg/jsonrpc"
)

// TODO: support longer IDs (or just use string as a key)
const maxIDLength = 32

type wsID [maxIDLength]byte

func transformID(id jsonrpc.ID) (wsID, error) {
	if len(id) > maxIDLength {
		return wsID{}, errors.New("id too long")
	}

	var wsID wsID
	copy(wsID[:], id)
	return wsID, nil
}
