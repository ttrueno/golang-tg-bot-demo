package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/x-goto/golang-tg-bot-demo/lib/e"
	"github.com/x-goto/golang-tg-bot-demo/storage"
)

const defaultPerm = 8776

type Storage struct {
	basePath string
}

func New(basePath string) *Storage {
	return &Storage{
		basePath: basePath,
	}
}

func (s *Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.Wrap("can't save page", err) }()

	fPath := filepath.Join(s.basePath, page.Username)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s *Storage) PickRandom(username string) (p *storage.Page, err error) {
	fPath := filepath.Join(s.basePath, username)

	files, err := os.ReadDir(fPath)
	if err != nil {
		return nil, err
	}

	filesNumber := len(files)

	if filesNumber == 0 {
		return nil, storage.ErrNoSaved
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(filesNumber)

	file := files[n]

	return s.decodePage(filepath.Join(fPath, file.Name()))
}

func (s *Storage) Remove(p *storage.Page) (err error) {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file: page hashing error", err)
	}

	fPath := filepath.Join(s.basePath, p.Username, fileName)
	if err := os.Remove(fPath); err != nil {
		msg := fmt.Sprintf("can't remove file: %s", fPath)
		return e.Wrap(msg, err)
	}

	return nil
}

func (s *Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file exists: page hashing error", err)
	}

	fPath := filepath.Join(s.basePath, p.Username, fileName)

	switch _, err = os.Stat(fPath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", fPath)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s *Storage) decodePage(filePath string) (st *storage.Page, err error) {
	defer func() { err = e.Wrap("can't decode file", err) }()

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, err
	}

	st = &p
	return st, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
