package auth

const (
	Authorization  = "Authorization"
	Bearer         = "Bearer"
	XJwtToken      = "X-JWT-Token"
	XAuthenticated = "X-Authenticated"
	XUser          = "X-User"
	XUserID        = "X-User-ID"
	XGroups        = "X-Groups"
	Role           = "Role"
	RoleAdmin      = "Admin"
	RoleEditor     = "Editor"
	RoleViewer     = "Viewer"
)

const (
	PermRead   = "read"
	PermWrite  = "write"
	PermUpdate = "update"
	PermDelete = "delete"
	PermList   = "list"
)

const (
	GroupGoogle = "google.com"
	GroupApple  = "apple.com"
)
