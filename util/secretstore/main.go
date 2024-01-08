package secretstore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/99designs/keyring"
	"github.com/rs/zerolog/log"
)

const (
	SecretStoreKey = "secret-store"
)

type SecretStore struct {
	ring keyring.Keyring

	shelf SecretShelf
}

type SecretShelf struct {
	key string

	isLoaded bool
	isDirty  bool

	Data []byte `json:"data"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (sh *SecretShelf) HasData() bool {
	return sh.Data != nil
}

func New(shelfKey string) (*SecretStore, error) {
	keyringPath, ok := os.LookupEnv("KEYRING_PATH")
	if !ok {
		ucd, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}

		keyringPath = filepath.Join(ucd, SecretStoreKey, "keyring")
	}

	ring, err := keyring.Open(keyring.Config{
		AllowedBackends: []keyring.BackendType{
			keyring.FileBackend,
		},

		ServiceName: SecretStoreKey,
		FileDir:     keyringPath,
		FilePasswordFunc: func(string) (string, error) {
			return SecretStoreKey, nil
		},
	})

	if err != nil {
		return nil, err
	}

	store := &SecretStore{
		ring: ring,

		shelf: SecretShelf{
			key:      shelfKey,
			isLoaded: false,
			isDirty:  false,
		},
	}

	return store, nil
}

func (s *SecretStore) GetShelfData() ([]byte, error) {
	if s.shelf.isLoaded {
		return s.shelf.Data, nil
	}

	storeItem, err := s.ring.Get(s.shelf.key)

	if err != nil {
		log.Warn().Err(err).Msg("no store item found")
		if err == keyring.ErrKeyNotFound {
			s.shelf.CreatedAt = time.Now()
			s.shelf.UpdatedAt = time.Now()
			s.shelf.isDirty = true
		} else {
			return s.shelf.Data, err
		}
	}

	if s.shelf.isDirty {
		// save shelf
		s.saveStore()
	} else {
		// load store item (onto) shelf
		err = json.Unmarshal(storeItem.Data, &s.shelf)
		if err != nil {
			return s.shelf.Data, err
		}
	}

	s.shelf.isLoaded = true
	s.shelf.isDirty = false

	return s.shelf.Data, err
}

func (s *SecretStore) UpdateShelf(shelfKey string, shelfData []byte) error {
	if s == nil {
		newStore, err := New(shelfKey)
		if err != nil {
			return err
		}

		s = newStore
	}

	s.shelf = SecretShelf{
		key:       s.shelf.key,
		Data:      shelfData,
		isDirty:   true,
		UpdatedAt: time.Now(),
	}

	err := s.saveStore()

	return err
}

func (s *SecretStore) saveStore() error {
	jsonBytes, err := json.Marshal(s.shelf)
	if err != nil {
		return err
	}

	newItem := keyring.Item{
		Key:         s.shelf.key,
		Data:        jsonBytes,
		Label:       fmt.Sprintf("%s-secret-shelf", SecretStoreKey),
		Description: fmt.Sprintf("a %s shelf for %s", SecretStoreKey, s.shelf.key),
	}

	err = s.ring.Set(newItem)

	if err != nil {
		return err
	}

	s.shelf.isDirty = false

	return nil
}
