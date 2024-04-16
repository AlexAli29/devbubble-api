package core

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ChatUser struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type PrivateChat struct {
	Id       string     `json:"id"`
	Users    []ChatUser `json:"users"`
	Messages []Message  `json:"messages"`
}

type CreatePrivateChatDto struct {
	SecondParticipantId string `json:"secondParticipantId" validate:"required,uuid"`
	UserId              string `json:"userId" validate:"required,uuid"`
}

type RemovePrivateChatDto struct {
	ChatId string `json:"chatId" validate:"required,uuid"`
	UserId string `json:"userId" validate:"required,uuid"`
}

type CreatePrivateChatResponse struct {
	Id string `json:"id" validate:"required,uuid"`
}

type GetPrivateChatsResponse struct {
	Id                    string           `json:"id"`
	Name                  string           `json:"name"`
	LastMessage           string           `json:"lastMessage"`
	LastMessageTime       pgtype.Timestamp `json:"lastMessageTime"`
	IsFromMe              bool             `json:"isFromMe"`
	ChatType              int              `json:"chatType"`
	LastMessageSenderName pgtype.Text      `json:"lastMessageSenderName"`
}

var (
	PRIVATE_CHAT = 1
	GROUP_CHAT   = 2
)
