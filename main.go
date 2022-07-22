/* IMPORTANT NOTE!!

When formatting text for entries, use _text_ for italic, ***text*** for bold, otherwise it won't work due to telegram API parser issues  */

package main

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/go-github/v45/github"
	//"strconv"
)

//TODO: Add more keyboards (for the god of keyboards)
var mainMenuKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Community"),
		tgbotapi.NewKeyboardButton("Fundrise"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Direct action"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Marketing"),
		tgbotapi.NewKeyboardButton("PR"),
	),
)

var PRKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("SMM"),
		tgbotapi.NewKeyboardButton("News management"),
	),
)

//to operate the bot, put a text file containing key for your bot acquired from telegram "botfather" to the same directory with this file
var tgApiKeyRaw, err = os.ReadFile(".secret")
var tgApiKeyString = bytes.NewBuffer(tgApiKeyRaw).String()
var tgApiKey = strings.Split(tgApiKeyString, "\n")

var bot, error1 = tgbotapi.NewBotAPI(string(tgApiKey[0]))

// session is a global struct for user session. it can be updated, loaded and saved independently to any user
type session struct {
	chat_id int64
	status  int64
}

// mapping from tgid to session
var userSession = make(map[int64]session)

func main() {

	//those are used to get info from our GitHub.
	ctx := context.Background()
	owner := "GreenLineProtest"
	repo := "Protester-b0t"
	client := github.NewClient(nil)

	var tgApiKey, _ = os.ReadFile(".secret")
	var bot, _ = tgbotapi.NewBotAPI(string(tgApiKey))

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// msgTemplates["hello"] = MsgTemplate{msg_string: "Hello, this is greating message"}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	//whenever bot gets a new message, check for user id in the database happens, if it's a new user, the entry in the database is created.
	for update := range updates {

		if update.Message != nil {

			//greetings part
			if _, ok := userSession[update.Message.From.ID]; !ok {

				userSession[update.Message.From.ID] = session{update.Message.Chat.ID, 0}
				path := "md/hello/README.md"
				readcontent, _, _, _ := client.Repositories.GetContents(ctx, owner, repo, path, nil)
				content := decode(readcontent)

				msg := tgbotapi.NewMessage(userSession[update.Message.From.ID].chat_id, content)
				//	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				msg.ReplyMarkup = mainMenuKeyboard
				msg.ParseMode = "Markdown"
				bot.Send(msg)

			} else {

				if userSession[update.Message.From.ID].status == 0 {

					//this block is used ONLY for switching keyboards, so categories system may be implemented.
					if update.Message.Text == "PR" ||
						// update.Message.Text == "Fundrise" ||
						// update.Message.Text == "Marketing" ||
						// update.Message.Text == "Community" ||
						update.Message.Text == "Direct action" {

						switch update.Message.Text {
						case "PR":
							msg := tgbotapi.NewMessage(userSession[update.Message.From.ID].chat_id, "Какое направление?")
							msg.ReplyMarkup = PRKeyboard
							bot.Request(msg)

						}

						//this block is used for sending actual text
					} else if update.Message.Text == "Community" ||
						update.Message.Text == "Fundrise" ||
						update.Message.Text == "Marketing" ||
						update.Message.Text == "Direct action" ||
						update.Message.Text == "News management" ||
						update.Message.Text == "SMM" {

						var content string

						switch update.Message.Text {
						case "Community":
							path := "md/community_managment/README.md"
							readcontent, _, _, _ := client.Repositories.GetContents(ctx, owner, repo, path, nil)
							content = decode(readcontent)

						case "Fundrise":
							path := "md/fundrise/README.md"
							readcontent, _, _, _ := client.Repositories.GetContents(ctx, owner, repo, path, nil)
							content = decode(readcontent)

						case "Marketing":
							path := "md/marketing/README.md"
							readcontent, _, _, _ := client.Repositories.GetContents(ctx, owner, repo, path, nil)
							content = decode(readcontent)

						case "Direct action":
							path := "md/bo/README.md"
							readcontent, _, _, _ := client.Repositories.GetContents(ctx, owner, repo, path, nil)
							content = decode(readcontent)

						//PR
						case "News management":
							path := "md/PR/NewsManagement.md"
							readcontent, _, _, _ := client.Repositories.GetContents(ctx, owner, repo, path, nil)
							content = decode(readcontent)

						case "SMM":
							path := "md/PR/SMM.md"
							readcontent, _, _, _ := client.Repositories.GetContents(ctx, owner, repo, path, nil)
							content = decode(readcontent)

						}

						//	content := readMd("test","test")
						msg := tgbotapi.NewMessage(userSession[update.Message.From.ID].chat_id, content)
						msg.ReplyMarkup = mainMenuKeyboard
						bot.Send(msg)

					}
				}

			}
		}
	}
}

//apply to client.Repositories.GetContents function output
func decode(repotext *github.RepositoryContent) (message string) {
	n := *repotext.Content
	m, _ := b64.StdEncoding.DecodeString(n)
	message = string(m)
	return message
}
