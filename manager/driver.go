package manager

// Driver is an interface for running maestro in
type Driver interface {
	Run() error
}
