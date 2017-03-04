package toml

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/tboerger/redirects/model"
	"github.com/tboerger/redirects/store"
	"io/ioutil"
	"time"
)

// redirectCollection represents the internal storage collection.
type redirectCollection struct {
	Redirects []*model.Redirect `toml:"redirects"`
}

// GetRedirects retrieves all redirects from the TOML store.
func (db *data) GetRedirects() ([]*model.Redirect, error) {
	root, err := db.load()

	if err != nil {
		return nil, err
	}

	return root.Redirects, nil
}

// GetRedirect retrieves a specific redirect from the TOML store.
func (db *data) GetRedirect(id string) (*model.Redirect, error) {
	root, err := db.load()

	if err != nil {
		return nil, err
	}

	for _, record := range root.Redirects {
		if record.ID == id {
			return record, nil
		}
	}

	return nil, store.ErrRedirectNotFound
}

// DeleteRedirect deletes a redirect from the TOML store.
func (db *data) DeleteRedirect(id string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	root, err := db.load()

	if err != nil {
		return err
	}

	for row, record := range root.Redirects {
		if record.ID == id {
			root.Redirects = append(
				root.Redirects[:row],
				root.Redirects[row+1:]...,
			)

			return db.write(root)
		}
	}

	return store.ErrRedirectNotFound
}

// UpdateRedirect updates a redirect on the TOML store.
func (db *data) UpdateRedirect(update *model.Redirect) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	root, err := db.load()

	if err != nil {
		return err
	}

	for row, record := range root.Redirects {
		if record.ID == update.ID {
			root.Redirects[row] = update
			return db.write(root)
		}
	}

	return store.ErrRedirectNotFound
}

// CreateRedirect creates a redirect on the TOML store.
func (db *data) CreateRedirect(create *model.Redirect) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	root, err := db.load()

	if err != nil {
		return err
	}

	for _, record := range root.Redirects {
		if record.Source == create.Source {
			return store.ErrRedirectSourceExists
		}
	}

	create.ID = fmt.Sprintf(
		"%x",
		md5.Sum([]byte(string(time.Now().Unix()))),
	)

	root.Redirects = append(
		root.Redirects,
		create,
	)

	return db.write(root)
}

// load parses all available records from the storage.
func (db *data) load() (*redirectCollection, error) {
	res := &redirectCollection{
		Redirects: make([]*model.Redirect, 0),
	}

	content, err := ioutil.ReadFile(db.file)

	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(string(content), res); err != nil {
		return nil, err
	}

	return res, nil
}

// write writes the TOML content back to the storage.
func (db *data) write(content *redirectCollection) error {
	buf := new(bytes.Buffer)

	if err := toml.NewEncoder(buf).Encode(content); err != nil {
		return err
	}

	if err := ioutil.WriteFile(db.file, buf.Bytes(), 0640); err != nil {
		return err
	}

	return nil
}