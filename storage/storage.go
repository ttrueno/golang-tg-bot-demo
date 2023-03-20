package storage

import (
	"crypto/sha256"
	"errors"
	"io"

	"github.com/x-goto/golang-tg-bot-demo/lib/e"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(username string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	Username string
}

var ErrNoSaved = errors.New("file's not saved")

func (p *Page) Hash() (hash string, err error) {
	defer func() { err = e.Wrap("can't calculate hash", err) }()
	h := sha256.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", err
	}

	if _, err := io.WriteString(h, p.Username); err != nil {
		return "", err
	}

	return string(h.Sum(nil)), nil
}
