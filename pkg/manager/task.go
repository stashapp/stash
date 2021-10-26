package manager

import "context"

type Task interface {
	Start(context.Context)
	GetDescription() string
}
