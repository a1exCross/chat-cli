package handler

import (
	"bufio"
	"context"
	"fmt"
	"github.com/a1exCross/chat-cli/internal/app"
	"github.com/a1exCross/chat-cli/internal/client/chat/proto/chat_v1"
	"github.com/a1exCross/chat-cli/internal/model"
	"github.com/a1exCross/chat-cli/internal/utils"
	"log"
	"os"
	"strings"
	"time"
)

type ChatHandler struct {
	di *app.ServiceProvider
}

func NewChatHandler(di *app.ServiceProvider) *ChatHandler {
	return &ChatHandler{di: di}
}

func (h *ChatHandler) Do(ctx context.Context, chatID int64, cfg model.UserInfoConfig) {
	// почему то здесь не возвращается ошибка...
	client, err := h.di.GetChatService().Connect(ctx, chatID, cfg.UserID)
	if err != nil {
		if utils.TokenHasExpired(cfg.AccessToken) {
			err = utils.TryReauthorize(h.di, cfg)
			if err != nil {
				log.Fatalf("failed to reauthorize: %v", err)
			}
		}
		log.Fatalf("falied to connect chat with id %d: %v", chatID, err)
	}

	go h.listenInputMessages(client)

	for {
		scanner := bufio.NewScanner(os.Stdin)
		var lines strings.Builder

		scanner.Scan()
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		lines.WriteString(line)
		lines.WriteString("\n")

		// удаление вводимой строки в консоль
		fmt.Printf("\033[1A\033[K")

		err = scanner.Err()
		if err != nil {
			log.Println("failed to scan message: ", err)
		}

		timeNow := time.Now()

		err = h.di.GetChatService().SendMessage(ctx, model.SendMessageParams{
			ChatID: chatID,
			Message: model.Message{
				Text:      lines.String(),
				From:      cfg.Username,
				CreatedAt: timeNow,
			},
		})

		// TODO: сюда в будущем прилетит проверка на случай, если токен будет не активен
		if err != nil {
			if utils.TokenHasExpired(cfg.AccessToken) {
				err = utils.TryReauthorize(h.di, cfg)
				if err != nil {
					log.Fatalf("failed to reauthorize: %v", err)
				}
			}
			log.Fatalf("failed to send message: %v", err)
		}

		fmt.Printf("%s at %v: %s", cfg.Username, h.getPrettyDateTime(timeNow.UTC()), lines.String())
	}
}

func (h *ChatHandler) listenInputMessages(client chat_v1.ChatV1_ConnectClient) {
	for {
		msg, err := client.Recv()
		if err != nil {
			log.Fatalf("failed to receive message: %v", err)
		}

		fmt.Printf("%s at %v: %s", msg.GetFrom(), h.getPrettyDateTime(msg.GetTimestamp().AsTime()), msg.GetText())
	}
}

func (h *ChatHandler) getPrettyDateTime(time time.Time) string {
	return time.Format("02-01 15:04")
}
