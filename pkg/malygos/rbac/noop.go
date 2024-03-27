package rbac

type NoopRBAC struct{}

func NewNoop() *NoopRBAC {
	return &NoopRBAC{}
}

func (n *NoopRBAC) IsAllowed(username, ressource, action string) bool {
	return true
}
