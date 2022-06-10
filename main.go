package main

import (
	"bytes"
	"log"
	"os"

	//"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//Those are keyboards that appear for some of the questions
var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("ERC20"),
		tgbotapi.NewKeyboardButton("ERC20Snapshot"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("ERC20Votes")),
)

var yesNoKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Yes"),
		tgbotapi.NewKeyboardButton("No")),
)

var correctKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Name")),

	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Symbol"),
		tgbotapi.NewKeyboardButton("Supply"),
		tgbotapi.NewKeyboardButton("Type")),

	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("It's all correct"),
	),
)

var mainMenuKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Community"),
		tgbotapi.NewKeyboardButton("Fundrise"),
		tgbotapi.NewKeyboardButton("Marketing"),
		tgbotapi.NewKeyboardButton("PR"),
		tgbotapi.NewKeyboardButton("DirectAction"),
	),
				
)


//to operate the bot, put a text file containing key for your bot acquired from telegram "botfather" to the same directory with this file
var tgApiKeyRaw, err = os.ReadFile(".secret")
var tgApiKeyString = bytes.NewBuffer(tgApiKeyRaw).String()
var tgApiKey = strings.Split(tgApiKeyString,"\n")

var bot, error1 = tgbotapi.NewBotAPI(string(tgApiKey[0]))


// session is a global struct for user session. it can be updated, loaded and saved independently to any user
type session struct {
	chat_id				int64
	status				int64

}

// mapping from tgid to session
var userSession = make(map[int64]session)

type MsgTemplate struct {
//	id                    int64
	msg_string			  string
}

var msgTemplates = make (map[string] MsgTemplate)


func readMd(dir string,name string) (result string) {
	base := "./md/"
	link := base + dir +"/" + name + ".md"
	page, err := os.ReadFile(link)
	if err != nil {
		log.Panic(err)
	}
	page_string := bytes.NewBuffer(page).String()
	return  page_string
}

func main() {

	bot, err = tgbotapi.NewBotAPI(string(tgApiKey[0]))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	msgTemplates["hello"] = MsgTemplate{msg_string: "Hello, this is greating message"}



	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	//whenever bot gets a new message, check for user id in the database happens, if it's a new user, the entry in the database is created.
	for update := range updates {

		if update.Message != nil {


			if update.Message.Text == "/start"{
				userSession[update.Message.From.ID] = session{update.Message.Chat.ID, 0}
				msg := tgbotapi.NewMessage(userSession[update.Message.From.ID].chat_id, msgTemplates["hello"].msg_string)
				msg.ReplyMarkup = mainMenuKeyboard
				bot.Send(msg)
			}

			if _, ok := userSession[update.Message.From.ID]; !ok {

			//	userDatabase[update.Message.From.ID] = user{update.Message.Chat.ID, 0, "", "", 0, 0, ""}
				userSession[update.Message.From.ID] = session{update.Message.Chat.ID, 0}
				msg := tgbotapi.NewMessage(userSession[update.Message.From.ID].chat_id, msgTemplates["hello"].msg_string)
			//	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				msg.ReplyMarkup = mainMenuKeyboard
				bot.Send(msg)
			} else {

				//first check for user status, (for a new user status 0 is set automatically), then user reply for the first bot message is logged to a database as name AND user status is updated
				if userSession[update.Message.From.ID].status == 0 {
					if update.Message.Text == "Community" || update.Message.Text == "Fundrise" || update.Message.Text == "Marketing" || update.Message.Text == "PR" || update.Message.Text == "DirectAction"  {
						var content string
						
						if update.Message.Text == "Community" {
							content = readMd("community_managment","README")
						} else if update.Message.Text == "Fundrise" {
							content = readMd("fundrise","README")
						} else if update.Message.Text == "Marketing" {
							content = readMd("marketing","README")
						} else if update.Message.Text == "PR" {
							content = readMd("PR","README")
						} else if update.Message.Text == "DirectAction" {
							content = readMd("bo","README")
						}
					
					
					
					if updateDb, ok := userSession[update.Message.From.ID]; ok {
					//	updateDb.exportTokenName = update.Message.Text
						updateDb.status = 1
						userSession[update.Message.From.ID] = updateDb
					}
				//	content := readMd("test","test")
					msg := tgbotapi.NewMessage(userSession[update.Message.From.ID].chat_id,content)
					msg.ReplyMarkup = mainMenuKeyboard
					bot.Send(msg)

				}

					//logic is that 1 incoming message fro the user equals one status check in database, so each status check ends with the message asking the next question
				} else if userSession[update.Message.From.ID].status == 1 {
				
					if update.Message.Text == "Community" || update.Message.Text == "Fundrise" || update.Message.Text == "Marketing" || update.Message.Text == "PR" || update.Message.Text == "DirectAction"  {
						var content string
						
						if update.Message.Text == "Community" {
							content = readMd("community_managment","README")
						} else if update.Message.Text == "Fundrise" {
							content = readMd("fundrise","README")
						} else if update.Message.Text == "Marketing" {
							content = readMd("marketing","README")
						} else if update.Message.Text == "PR" {
							content = readMd("PR","README")
						} else if update.Message.Text == "DirectAction" {
							content = readMd("bo","README")
						}
						
				
					if updateDb, ok := userSession[update.Message.From.ID]; ok {
					//	updateDb.exportTokenSymbol = update.Message.Text
						updateDb.status = 0
						userSession[update.Message.From.ID] = updateDb
					}
					
					msg := tgbotapi.NewMessage(userSession[update.Message.From.ID].chat_id, content)
					msg.ReplyMarkup = mainMenuKeyboard
					bot.Send(msg)
																									}
				}
			}
	}	}
}
