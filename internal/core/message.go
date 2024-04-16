package core

type Message struct {
	Id         string `json:"id"`
	Text       string `json:"text"`
	SenderId   string `json:"senderId"`
	SenderName string `json:"senderName"`
	IsFromMe   bool   `json:"isFromMe"`
	ChatId     string `json:"chatId"`
	CreatedAt  string `json:"createdAt"`
}

type CreateMessageDto struct {
	Text   string `json:"text" validate:"required"`
	ChatId string `json:"chatId" validate:"required,uuid"`
	UserId string `json:"userId" validate:"required,uuid"`
}

type RemoveMessageDto struct {
	Id     string `json:"id" validate:"required,uuid"`
	UserId string `json:"userId" validate:"required,uuid"`
}

type GetMessagesByChatIdDto struct {
	ChatId string `json:"chatId" validate:"required,uuid"`
	UserId string `json:"userId" validate:"required,uuid"`
}
