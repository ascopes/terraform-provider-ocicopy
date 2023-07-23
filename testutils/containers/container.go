package containers

type Container interface {
	AddExposedPort(port string)
}
