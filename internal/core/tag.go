package core

type Tag struct {
	Id   string `json:"id"`
	Icon string `json:"icon"`
	Text string `json:"text"`
}
type AddTagDto struct {
	UserId string `json:"userId" validate:"required,uuid"`
	TagId  string `json:"tagId" validate:"required,uuid"`
}
type RemoveTagDto struct {
	UserId string `json:"userId" validate:"required,uuid"`
	TagId  string `json:"tagId" validate:"required,uuid"`
}
