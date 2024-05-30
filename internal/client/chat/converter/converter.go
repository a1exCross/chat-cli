package converter

import (
	"github.com/a1exCross/chat-cli/internal/client/chat/proto/chat_v1"
	"github.com/a1exCross/chat-cli/internal/model"
)

func ProtoToChats(chatPB []*chat_v1.Chat) []model.Chat {
	var chats = make([]model.Chat, len(chatPB))
	for i := 0; i < len(chatPB); i++ {
		chats[i] = model.Chat{
			Usernames: chatPB[i].Usernames,
			ID:        chatPB[i].Id,
		}
	}

	return chats
}
