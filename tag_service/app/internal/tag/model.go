package tag

type Tag struct {
	ID      int    `json:"id" bson:"_id,omitempty"`
	Name    string `json:"name" bson:"name,omitempty"`
	Color   string `json:"color" bson:"color,omitempty"`
	OwnerID string `json:"owner_id" bson:"owner_id,omitempty"`
}

func NewTag(dto CreateTagDTO) Tag {
	return Tag{
		Name:    dto.Name,
		Color:   dto.Color,
		OwnerID: dto.OwnerID,
	}
}

func UpdatedTag(dto UpdateTagDTO) Tag {
	return Tag{
		ID:      dto.ID,
		Name:    dto.Name,
		Color:   dto.Color,
	}
}

type CreateTagDTO struct {
	Name    string `json:"name" bson:"name"`
	Color   string `json:"color" bson:"color"`
	OwnerID string `json:"owner_id" bson:"owner_id"`
}

type UpdateTagDTO struct {
	ID    int    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
	Color string `json:"color,omitempty" bson:"color,omitempty"`
}
