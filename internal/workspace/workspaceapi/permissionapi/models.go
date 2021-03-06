package permissionapi

type CreatePermsRequest struct {
	UserIDs []uint64 `form:"user_ids" json:"user_ids" binding:"required"`
}

type GetPermsRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

type DeletePermRequest struct {
	UserID uint64 `json:"user_id" form:"user_id" binding:"required"`
	Type   int    `json:"type" form:"type" binding:"required"`
}

type PermissionResponse struct {
	CreatedAt int64  `json:"created_at"`
	UserID    uint64 `json:"user_id"`
	Role      int32  `json:"role"`
}

type GetPermissionResponse struct {
	Total   int                  `json:"total"`
	Members []PermissionResponse `json:"members"`
}
