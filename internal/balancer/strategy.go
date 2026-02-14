package balancer

type Strategy interface{
	NewBackend() *Backend
}