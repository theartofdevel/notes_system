package category_service

type CreateCategoryDTO struct {
	Name       string `json:"name"`
	UserUuid   string `json:"user_uuid"`
	ParentUuid string `json:"parent_uuid,omitempty"`
}

type UpdateCategoryDTO struct {
	Uuid       string `json:"uuid"`
	Name       string `json:"name"`
	UserUuid   string `json:"user_uuid"`
	ParentUuid string `json:"parent_uuid,omitempty"`
}

type DeleteCategoryDTO struct {
	Uuid     string `json:"uuid"`
	UserUuid string `json:"user_uuid"`
}
