package note

type Note struct {
	UUID         string `json:"uuid" bson:"_id,omitempty"`
	Header       string `json:"header" bson:"header,omitempty"`
	Body         string `json:"body,omitempty" bson:"body,omitempty"`
	ShortBody    string `json:"short_body,omitempty" bson:"short_body,omitempty"`
	CategoryUUID string `json:"category_uuid" bson:"category_uuid,omitempty"`
	Tags         []int  `json:"tags" bson:"tags,omitempty"`
}

func (cn *Note) GenerateShortBody() {
	var shortLen int
	if len(cn.Body) > 1000 {
		shortLen = 300
	} else {
		shortLen = len(cn.Body)
	}
	cn.ShortBody = cn.Body[0:shortLen]
}

func NewNote(dto CreateNoteDTO) Note {
	return Note{
		Header:       dto.Header,
		Body:         dto.Body,
		CategoryUUID: dto.CategoryUUID,
		Tags:         dto.Tags,
	}
}

func UpdatedNote(dto UpdateNoteDTO) Note {
	return Note{
		UUID:         dto.UUID,
		Header:       dto.Header,
		Body:         dto.Body,
		CategoryUUID: dto.CategoryUUID,
		Tags:         dto.Tags,
	}
}

type CreateNoteDTO struct {
	Header       string `json:"header" bson:"header"`
	Body         string `json:"body" bson:"body"`
	CategoryUUID string `json:"category_uuid" bson:"category_uuid"`
	Tags         []int  `json:"tags" bson:"tags"`
}

type UpdateNoteDTO struct {
	UUID         string `json:"uuid" bson:"_id,omitempty"`
	Header       string `json:"header,omitempty" bson:"header,omitempty"`
	Body         string `json:"body,omitempty" bson:"body,omitempty"`
	CategoryUUID string `json:"category_uuid,omitempty" bson:"category_uuid,omitempty"`
	Tags         []int  `json:"tags,omitempty" bson:"tags,omitempty"`
}
