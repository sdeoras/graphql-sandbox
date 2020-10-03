package auth

func Permissions(permissions ...string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, permission := range permissions {
		m[permission] = struct{}{}
	}

	return m
}
