## 1.Access Control List

```go
package acl

import (
"errors"
"strings"
)

// Role represents a user role
type Role string

const (
RoleAdmin    Role = "admin"
RoleEditor   Role = "editor"
RoleViewer   Role = "viewer"
RoleCustomer Role = "customer"
)

// Permission represents an action on a resource
type Permission struct {
Resource string
Action   string
}

// ACL stores role-based permissions
type ACL struct {
permissions map[Role][]Permission
}

// NewACL creates a new ACL with default permissions
func NewACL() *ACL {
acl := &ACL{
permissions: make(map[Role][]Permission),
}

	// Define permissions for each role
	acl.permissions[RoleAdmin] = []Permission{
		{"*", "*"}, // Admin can do anything
	}

	acl.permissions[RoleEditor] = []Permission{
		{"articles", "create"},
		{"articles", "edit"},
		{"articles", "view"},
	}

	acl.permissions[RoleViewer] = []Permission{
		{"articles", "view"},
	}

	acl.permissions[RoleCustomer] = []Permission{
		{"orders", "create"},
		{"orders", "view"},
		{"profile", "view"},
		{"profile", "edit"},
	}

	return acl
}

// HasPermission checks if a role has permission to perform an action on a resource
func (a *ACL) HasPermission(role Role, resource, action string) bool {
// Check if role exists
permissions, ok := a.permissions[role]
if !ok {
return false
}

	// Check for wildcard permissions (admin case)
	for _, p := range permissions {
		if (p.Resource == "*" || p.Resource == resource) && 
		   (p.Action == "*" || p.Action == action) {
			return true
		}
	}

	return false
}

// GetRoleFromScope converts OAuth2 scope to ACL role
func GetRoleFromScope(scope string) (Role, error) {
scopes := strings.Split(scope, " ")
for _, s := range scopes {
switch s {
case "admin":
return RoleAdmin, nil
case "editor":
return RoleEditor, nil
case "viewer":
return RoleViewer, nil
case "customer":
return RoleCustomer, nil
}
}
return "", errors.New("no valid role found in scope")
}
```
## 2.Middleware for Authorization
