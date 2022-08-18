package uuid

import (
	uuid1 "github.com/google/uuid"
)

var (
	InvalidUUID    = uuid1.UUID{}
	InvalidUUIDStr = InvalidUUID.String()
)
