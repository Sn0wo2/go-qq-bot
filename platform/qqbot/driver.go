package runtime

type Driver interface {
	Run() error
	Stop() error
}
