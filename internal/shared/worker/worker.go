package worker

type IWorker[T any] interface {
	Run(poolCount int)
	AddTask(value T)
	Close()
}
