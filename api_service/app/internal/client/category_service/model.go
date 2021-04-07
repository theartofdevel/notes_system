package category_service

type CreateCategoryDTO struct {
	Name       string `json:"name"`
	UserUuid   string `json:"user_uuid"`
	ParentUuid string `json:"parent_uuid,omitempty"`
}

type UpdateCategoryDTO struct {
	Uuid       string `json:"uuid,omitempty"`
	Name       string `json:"name,omitempty"`
	UserUuid   string `json:"user_uuid,omitempty"`
	ParentUuid string `json:"parent_uuid,omitempty"`
}

type DeleteCategoryDTO struct {
	Uuid     string `json:"uuid"`
	UserUuid string `json:"user_uuid"`
}
