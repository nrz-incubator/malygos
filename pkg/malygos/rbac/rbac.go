package rbac

type RBAC interface {
	IsAllowed(username, ressource, action string) bool
}
