package root

import (
	"context"
	"fmt"
	"github.com/a1exCross/chat-cli/internal/app"
	"github.com/a1exCross/chat-cli/internal/config"
	"github.com/a1exCross/chat-cli/internal/handler"
	"github.com/a1exCross/chat-cli/internal/model"
	"github.com/a1exCross/chat-cli/internal/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "chat-cli",
	Short: "CLI утилита для работы с приложением",
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "что-то создает",
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "что-то удаляет",
}

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Создает нового пользователя",
	Run: func(cmd *cobra.Command, args []string) {
		usernamesStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatal("failed to get username")
		}

		/*serviceProvider := app.NewServiceProvider()

		serviceProvider.GetAuthService().CreateUser(context.Background(), model.UserCreateParams{})*/
		log.Printf(`username "%s" created`, usernamesStr)
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Удаляет пользователя",
	Run: func(cmd *cobra.Command, args []string) {
		userID, err := cmd.Flags().GetInt64("user_id")
		if err != nil {
			log.Fatal("failed to get user_id")
		}

		serviceProvider := app.NewServiceProvider()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		ctx = utils.NewOutgoingContextWithToken(context.Background(), config.LoadExecConfig().AccessToken)

		err = serviceProvider.GetAuthService().DeleteUser(ctx, userID)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf(`username "%d" deleted`, userID)
	},
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Подключиться к чату",
	Run: func(cmd *cobra.Command, args []string) {
		chatID, err := cmd.Flags().GetInt64("chat_id")
		if err != nil {
			log.Fatal("failed to get chat_id")
		}

		cfg := config.LoadExecConfig()

		serviceProvider := app.NewServiceProvider()

		ctx := utils.NewOutgoingContextWithToken(context.Background(), cfg.AccessToken)

		handler.NewChatHandler(serviceProvider).Do(ctx, chatID, cfg)
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Авторизация",
	Run: func(cmd *cobra.Command, args []string) {
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("failed to get login")
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Fatalf("failed to get login")
		}

		servProvider := app.NewServiceProvider()
		authService := servProvider.GetAuthService()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		refreshToken, err := authService.Login(ctx, model.UserInfo{
			Username: username,
			Password: password,
		})
		if err != nil {
			log.Fatalf("failed to login: %v", err)
		}

		accessToken, err := authService.Authorize(ctx, refreshToken)
		if err != nil {
			log.Fatalf("failed to authorize: %v", err)
		}

		userID, err := utils.GetUserIDFromAccessToken(accessToken)
		if err != nil {
			log.Fatalf("failed to get user id: %v", err)
		}

		userData := model.UserInfoConfig{
			UserID:       userID,
			AccessToken:  accessToken,
			Username:     username,
			RefreshToken: refreshToken,
		}

		config.SaveExecConfig(userData)

		fmt.Printf(color.GreenString("\nHello, %s!\n\n", username))
	},
}

var createChat = &cobra.Command{
	Use:   "chat",
	Short: "Создание чата",
	Run: func(cmd *cobra.Command, args []string) {
		usernames, err := cmd.Flags().GetStringArray("usernames")
		if err != nil {
			log.Fatal("failed to get usernames")
		}

		servProvider := app.NewServiceProvider()
		chatService := servProvider.GetChatService()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		ctx = utils.NewOutgoingContextWithToken(context.Background(), config.LoadExecConfig().AccessToken)

		chatID, err := chatService.Create(ctx, usernames)
		if err != nil {
			log.Fatalf("failed to create chat: %v", err)
		}

		cmd.Printf("chat created with id %d\n", chatID)
	},
}

var deleteChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Удаление чата",
	Run: func(cmd *cobra.Command, args []string) {
		chatID, err := cmd.Flags().GetInt64("chat_id")
		if err != nil {
			log.Fatal("failed to get chat_id")
		}

		servProvider := app.NewServiceProvider()
		chatService := servProvider.GetChatService()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		ctx = utils.NewOutgoingContextWithToken(context.Background(), config.LoadExecConfig().AccessToken)

		err = chatService.Delete(ctx, chatID)
		if err != nil {
			log.Fatalf("failed to delete chat: %v", err)
		}

		cmd.Printf("chat deleted with id %d\n", chatID)
	},
}

var listChatsCMD = &cobra.Command{
	Use:   "list",
	Short: "Мои чаты",
	Run: func(cmd *cobra.Command, args []string) {
		servProvider := app.NewServiceProvider()

		// надо думать над оберткой, везде не хочется прописывать таймаут
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cfg := config.LoadExecConfig()

		ctx = utils.NewOutgoingContextWithToken(context.Background(), cfg.AccessToken)

		chats, err := servProvider.GetChatService().ListChats(ctx, cfg.Username)
		if err != nil {
			log.Fatal(err)
		}

		if len(chats) == 0 {
			fmt.Println("chats not found")

			return
		}

		for i := 0; i < len(chats); i++ {
			cmd.Printf("chatID: %d \t usernames: %v\n", chats[i].ID, chats[i].Usernames)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(listChatsCMD)

	initCreateUserCMD()
	initDeleteUserCMD()
	initCreateChatCMD()
	initDeleteChatCMD()
	initConnectCMD()
	initLoginCMD()
}

func initLoginCMD() {
	loginCmd.Flags().StringP("username", "u", "", "Логин")
	loginCmd.Flags().StringP("password", "p", "", "Пароль")

	err := loginCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark required flag: %s", err)
	}

	err = loginCmd.MarkFlagRequired("password")
	if err != nil {
		log.Fatalf("failed to mark required flag: %s", err)
	}
}

func initConnectCMD() {
	connectCmd.Flags().Int64("chat_id", 0, "Идентификатор чата")

	err := connectCmd.MarkFlagRequired("chat_id")
	if err != nil {
		log.Fatalf("failed to mark required flag: %s", err)
	}
}

func initDeleteChatCMD() {
	deleteCmd.AddCommand(deleteChatCmd)

	deleteChatCmd.Flags().Int64("chat_id", 0, "Идентификатор чата")
	err := deleteChatCmd.MarkFlagRequired("chat_id")
	if err != nil {
		log.Fatalf("failed to mark required flag: %s", err)
	}
}

func initDeleteUserCMD() {
	deleteCmd.AddCommand(deleteUserCmd)

	deleteUserCmd.Flags().Int64("user_id", 0, "Идентификатор пользователя")
	err := deleteUserCmd.MarkFlagRequired("user_id")
	if err != nil {
		log.Fatalf("failed to mark required flag: %s", err)
	}
}

func initCreateUserCMD() {
	createCmd.AddCommand(createUserCmd)

	createUserCmd.Flags().StringP("username", "u", "", "Логин пользователя")
	err := createUserCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark required flag: %s", err)
	}

	createUserCmd.Flags().StringP("name", "n", "", "Имя пользователя")
	err = createUserCmd.MarkFlagRequired("name")
	if err != nil {
		log.Fatalf("failed to mark required flag: %s", err)
	}

	createUserCmd.Flags().StringP("email", "e", "", "Электронная почта")
	err = createUserCmd.MarkFlagRequired("email")
	if err != nil {
		log.Fatalf("failed to mark required flag: %s", err)
	}

	createUserCmd.Flags().StringP("password", "p", "", "Пароль")
	err = createUserCmd.MarkFlagRequired("password")
	if err != nil {
		log.Fatalf("failed to mark required flag: %s", err)
	}

	createUserCmd.Flags().StringP("confirm", "c", "", "Подтверждение пароля")
	err = createUserCmd.MarkFlagRequired("confirm")
	if err != nil {
		log.Fatalf("failed to mark required flag: %s", err)
	}
}

func initCreateChatCMD() {
	createCmd.AddCommand(createChat)

	createChat.Flags().StringArray("usernames", []string{}, "Участники чата")
	err := createChat.MarkFlagRequired("usernames")
	if err != nil {
		log.Fatalf("failed to mark required flag: %s", err)
	}
}
