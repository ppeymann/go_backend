package utils

import "errors"

const ContextUserKey = "CONTEXT_USER"
const ContextRoleKey = "CONTEXT_ROLE"

var ErrUserPrincipalsNotFound = errors.New("UserPrincipals not found in context")
