package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go/jetstream"
)

var (
	ErrNotFound = errors.New("not found")
)

type Item struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (i *Item) Validate() {

}

type ItemStore struct {
	db    *pgxpool.Pool
	cache jetstream.KeyValue
}

func NewItemStore(db *pgxpool.Pool, cache jetstream.KeyValue) *ItemStore {
	return &ItemStore{
		db:    db,
		cache: cache,
	}
}

// TODO: Maybe implement event sourcing as the canonical store, but for now - just be lazy

func (s *ItemStore) Create(ctx context.Context, item *Item) error {
	if err := s.db.QueryRow(ctx, "INSERT INTO items (name) VALUES ($1) RETURNING id", item.Name).Scan(&item.ID); err != nil {
		return err
	}

	// Add cache entry
	itemJson, err := json.Marshal(item)
	if err != nil {
		return err
	}

	if _, err = s.cache.Create(ctx, item.ID, itemJson); err != nil {
		return err
	}

	return nil
}

func (s *ItemStore) Get(ctx context.Context, id string) (*Item, error) {
	cItem, err := s.cache.Get(ctx, id)
	if err != nil {
		if !errors.Is(err, jetstream.ErrKeyNotFound) || !errors.Is(err, jetstream.ErrKeyDeleted) {
			return nil, err
		}

		// Cache miss, head to the database
		var item Item
		if err = s.db.QueryRow(ctx, "SELECT id, name FROM items WHERE id = $1", id).
			Scan(&item.ID, &item.Name); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
		}

		return &item, nil
	}

	var item *Item
	if err = json.Unmarshal(cItem.Value(), item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *ItemStore) Update(ctx context.Context, id string, item *Item) error {
	if err := s.db.QueryRow(ctx, "UPDATE items SET name = $1 WHERE id = $2 RETURNING id", item.Name, id).Scan(&item.ID); err != nil {
		return err
	}

	// Update cache
	itemJson, err := json.Marshal(item)
	if err != nil {
		return err
	}

	if _, err = s.cache.Put(ctx, id, itemJson); err != nil {
		return err
	}

	return nil
}
