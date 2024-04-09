package main

import (
	"bakalover/hikari-bot/dict/jisho"
	"bakalover/hikari-bot/game"
	"fmt"
	"log"
	"os"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	HelpInfo        = "Справка по коммандам:\n/help - Справка по командам\n/sh_start - Начать игру\n/sh_stop - Закончить игру и вывести результаты"
	Unknown         = "Неизвестная команда"
	ShiritoryPrefix = "sh_"
)

func HandleCommand(dbConn *gorm.DB, bot *tg.BotAPI, msg *tg.Message) {
	command := msg.Command()

	if strings.HasPrefix(command, ShiritoryPrefix) {
		game.RunGameCommand(game.MsgContext{DbConn: dbConn, Bot: bot, Msg: msg})
		return
	}

	//Filter non-game commands e.g /help
	switch command {
	case "help":
		bot.Send(tg.NewMessage(msg.Chat.ID, HelpInfo))

	default:
		bot.Send(tg.NewMessage(msg.Chat.ID, Unknown))
	}
}

func main() {

	bot, err := tg.NewBotAPI(os.Getenv("HIKARI_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Couldn't initialize bot api!\n%v", err)
	}

	dsn := fmt.Sprintf("host=localhost user=%v password=%v dbname=%v port=5432 sslmode=disable", os.Getenv("PG_LOGIN"), os.Getenv("PG_LOGIN"), os.Getenv("PG_DB"))
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Couldn't establish connection to PostgreSQL!\n%v", err)
	}

	// bot.Debug = true

	uCfg := tg.NewUpdate(0) // No timeout (or maybe specify later)
	uCfg.Timeout = 60
	uCfg.AllowedUpdates = []string{"message"}


	// Strand | MPSC
	upds := bot.GetUpdatesChan(uCfg)

	dict := &jisho.JishoDict{}

	for upd := range upds {
		if msg := upd.Message; msg != nil {
			log.Printf("User: %v, Message: %v", msg.From.UserName, msg.Text)
			if msg.IsCommand() {
				HandleCommand(dbConn, bot, msg)
			} else {
				log.Println(game.Chat())
				log.Println(msg.Chat.ID)
				if game.Chat() == msg.Chat.ID && game.IsRunning() {
					game.HandleNextWord(game.MsgContext{DbConn: dbConn, Bot: bot, Msg: msg}, dict)
				}
			}
		}
	}
}
