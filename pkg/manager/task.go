package manager

type Task interface {
	Start()
	GetDescription() string
}
