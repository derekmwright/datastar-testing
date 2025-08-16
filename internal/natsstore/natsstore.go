package natsstore

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsStore struct {
	client jetstream.KeyValue
	prefix string
	logger *slog.Logger
}

type Option func(*NatsStore)

func WithPrefix(prefix string) func(*NatsStore) {
	return func(store *NatsStore) {
		store.prefix = prefix
	}
}

func New(ctx context.Context, js jetstream.JetStream, cfg jetstream.KeyValueConfig, opts ...Option) (*NatsStore, error) {
	var (
		kv  jetstream.KeyValue
		err error
	)

	if kv, err = js.CreateKeyValue(ctx, cfg); err != nil {
		return nil, err
	}

	store := &NatsStore{
		client: kv,
		prefix: "scs.session.",
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	for _, opt := range opts {
		opt(store)
	}

	return store, nil
}

func Must(store *NatsStore, err error) *NatsStore {
	if err != nil {
		panic(err)
	}

	return store
}

// CommitCtx adds a session token and data to the NATS KV store.
// Note: the TTL param is discarded as TTL is a bucket property in NATS.
func (s *NatsStore) CommitCtx(ctx context.Context, token string, b []byte, _ time.Time) error {
	if _, err := s.client.Create(ctx, s.prefix+token, b); err != nil {
		if errors.Is(err, nats.ErrKeyExists) {
			_, err = s.client.Put(ctx, s.prefix+token, b)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

// DeleteCtx removes a token and data from the NATS KV store.
func (s *NatsStore) DeleteCtx(ctx context.Context, token string) error {
	if err := s.client.Purge(ctx, s.prefix+token); err != nil {
		if errors.Is(err, nats.ErrKeyNotFound) {
			return nil
		}
		return err
	}
	return nil
}

// FindCtx finds a token and data from the NATS KV store.
func (s *NatsStore) FindCtx(ctx context.Context, token string) ([]byte, bool, error) {
	var (
		entry jetstream.KeyValueEntry
		err   error
	)

	if entry, err = s.client.Get(ctx, s.prefix+token); err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) || errors.Is(err, jetstream.ErrKeyDeleted) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return entry.Value(), true, nil
}

// Required for interface conformance

func (s *NatsStore) Delete(token string) (err error) {
	panic("Delete called, use DeleteCtx instead")
}

func (s *NatsStore) Find(token string) (b []byte, found bool, err error) {
	panic("Find called, use FindCtx instead")
}

func (s *NatsStore) Commit(token string, b []byte, expiry time.Time) (err error) {
	panic("Commit called, use CommitCtx instead")
}
