package v1dto

type CreateUserRequest struct {
	UserName string `json:"userName" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	UserName *string `json:"userName" binding:"omitempty,required"`
	Name     *string `json:"name" binding:"omitempty,required"`
	Email    string  `json:"email" binding:"required,email"`
	Password *string `json:"password" binding:"omitempty,required,min=6"`
}

type GetUserIdParams struct {
	UUID string `uri:"uuid" binding:"required,uuid"`
}

type ListUsersQuery struct {
	Page   int     `form:"page" binding:"omitempty,min=1"`
	Limit  int     `form:"limit" binding:"omitempty,min=1,max=100"`
	SortBy string  `form:"sort_by" binding:"omitempty,oneof=created_at updated_at name email"`
	Order  string  `form:"order" binding:"omitempty,oneof=asc desc"`
	Search *string `form:"search" binding:"omitempty,max=100"`
}

func (q *ListUsersQuery) Normalize() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Order == "" {
		q.Order = "asc"
	}
	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
}
