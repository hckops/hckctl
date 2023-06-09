package box

import (
	"context"
)

type Connection struct {
	ctx context.Context
	Out chan string // TODO
}

// TODO move client implementations here? + newBoxClient factory?

type BoxClient interface {
	Open() (*Connection, error)
}
