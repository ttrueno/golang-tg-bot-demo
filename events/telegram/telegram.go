package telegram

import (
	"errors"

	"github.com/x-goto/golang-tg-bot-demo/clients/telegram"
	"github.com/x-goto/golang-tg-bot-demo/events"
	"github.com/x-goto/golang-tg-bot-demo/lib/e"
	"github.com/x-goto/golang-tg-bot-demo/storage"
)

type EventHandler struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnkownEventType = errors.New("unkown event type")
	ErrUnkownMetaType  = errors.New("unkown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *EventHandler {
	return &EventHandler{
		tg:      client,
		storage: storage,
	}
}

func (eh *EventHandler) Fetch(limit int) ([]events.Event, error) {
	upds, err := eh.tg.Updates(eh.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	lenUpds := len(upds)

	if lenUpds == 0 {
		return nil, nil
	}

	evs := make([]events.Event, 0, len(upds))
	for _, upd := range upds {
		evs = append(evs, event(upd))
	}

	eh.offset = upds[lenUpds-1].ID + 1

	return evs, nil
}

func (eh *EventHandler) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return eh.processMessage(event)
	default:
		return e.Wrap("can't process event", ErrUnkownEventType)
	}
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	e := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		e.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return e
}

func (eh *EventHandler) processMessage(event events.Event) (err error) {
	defer func() { err = e.Wrap("can't process message", err) }()

	meta, err := meta(event)
	if err != nil {
		return err
	}

	if err := eh.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return err
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, ErrUnkownEventType
	}

	return res, nil
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}
