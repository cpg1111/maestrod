package manager

// Driver is an interface for running maestro in
type Driver interface {
	Run(name, confTarget, hostVolume string, args []string) error
	DestroyWorker(project, branch string) error
}
