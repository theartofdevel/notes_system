package note

type Note struct {
	UUID         string `json:"uuid" bson:"_id"`
	Header       string `json:"header" bson:"header"`
	Body         string `json:"body,omitempty" bson:"body"`
	ShortBody    string `json:"short_body,omitempty" bson:"short_body"`
	CategoryUUID string `json:"category_uuid" bson:"category_uuid"`
}

type CreateNoteDTO struct {
	Header       string `json:"header" bson:"header"`
	Body         string `json:"body" bson:"body"`
	ShortBody    string `json:"short_body,omitempty" bson:"short_body,omitempty"`
	CategoryUUID string `json:"category_uuid" bson:"category_uuid"`
}

func (cn *CreateNoteDTO) GenerateShortBody() {
	var shortLen int
	if len(cn.Body) > 1000 {
		shortLen = 300
	} else {
		shortLen = len(cn.Body)
	}
	cn.ShortBody = cn.Body[0:shortLen]
}

type UpdateNoteDTO struct {
	Header       string `json:"header,omitempty" bson:"header,omitempty"`
	Body         string `json:"body,omitempty" bson:"body,omitempty"`
	CategoryUUID string `json:"category_uuid,omitempty" bson:"category_uuid,omitempty"`
}
