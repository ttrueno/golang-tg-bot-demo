package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/x-goto/golang-tg-bot-demo/clients/telegram"
	"github.com/x-goto/golang-tg-bot-demo/lib/e"
	"github.com/x-goto/golang-tg-bot-demo/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (eh *EventHandler) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	if isAddCmd(text) {
		return eh.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return eh.sendRandom(chatID, username)
	case HelpCmd:
		return eh.sendHelp(chatID)
	case StartCmd:
		return eh.sendHello(chatID)
	default:
		return eh.tg.SendMessage(chatID, msgUnkownCommand)
	}
}

func (eh *EventHandler) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.Wrap("can't do command: save page", err) }()

	send := NewMessageSender(chatID, eh.tg)

	page := &storage.Page{
		URL:      pageURL,
		Username: username,
	}

	isExist, err := eh.storage.IsExists(page)
	if err != nil {
		return err
	}

	if isExist {
		return send(msgAlreadyExists)
	}

	if err := eh.storage.Save(page); err != nil {
		return err
	}

	if err := send(msgSaved); err != nil {
		return err
	}

	return nil
}

func (eh *EventHandler) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.Wrap("can't do command: send random", err) }()

	page, err := eh.storage.PickRandom(username)

	if err != nil && !errors.Is(err, storage.ErrNoSaved) {
		log.Print("hehe")
		return err
	}
	if errors.Is(err, storage.ErrNoSaved) {
		return eh.tg.SendMessage(chatID, msgNoSavedPage)
	}

	if err := eh.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return eh.storage.Remove(page)
}

func (eh *EventHandler) sendHello(chatID int) error {
	return eh.tg.SendMessage(chatID, msgHello)
}

func (eh *EventHandler) sendHelp(chatID int) error {
	return eh.tg.SendMessage(chatID, msgHelp)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}

func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}
