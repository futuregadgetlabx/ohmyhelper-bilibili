package task

type Task interface {
	Run()
	Name() string
}
