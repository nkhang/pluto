package errors

const (
	WorkspaceNotFound ErrorType = -(1300 + iota)
	WorkspaceQueryError
	WorkspaceErrorCreating
	WorkspaceCannotUpdate
	WorkspacePermissionErrorCreating
	WorkspacePermissionNotFound
	WorkspaceErrorDeleting
	WorkspacePermissionDeletingError
)
