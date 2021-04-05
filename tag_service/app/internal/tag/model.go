package tag

type Tag struct {
	ID       int    `json:"_id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Color    string `json:"color" bson:"color"`
	UserUUID string `json:"user_uuid" bson:"user_uuid"`
}

type CreateTagDTO struct {
	ID       int    `json:"_id,omitempty" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Color    string `json:"color" bson:"color"`
	UserUUID string `json:"user_uuid" bson:"user_uuid"`
}

type UpdateTagDTO struct {
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Color    string `json:"color,omitempty" bson:"color,omitempty"`
}

