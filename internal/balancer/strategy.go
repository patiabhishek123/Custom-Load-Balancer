package balancer

type Strategy interface {
	NextBackend() *Backend
	Pool() *BackendPool
}
