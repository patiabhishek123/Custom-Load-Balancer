package balancer

type strategy interface{
	NewBackend() *Backend
}