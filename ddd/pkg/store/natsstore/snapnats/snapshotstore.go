package snapnats

import (
	"context"
	"ddd/pkg/store"
	"ddd/pkg/store/natsstore"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

type snapshotStore[T any] struct {
	tname      string
	boundedCtx string
	kv         jetstream.KeyValue
}

func NewSnapshotStore[T any](ctx context.Context, js jetstream.JetStream) *snapshotStore[T] {
	aname, bname := natsstore.MetaFromType[T]()
	store := &snapshotStore[T]{
		tname:      aname,
		boundedCtx: bname,
	}
	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:  store.snapshotBucketName(),
		Storage: jetstream.MemoryStorage,
	})
	if err != nil {
		panic(err)
	}

	store.kv = kv
	return store
}

func (s *snapshotStore[T]) snapshotBucketName() string {
	return fmt.Sprintf("snapshot-%s-%s", s.boundedCtx, s.tname)
}

func (s *snapshotStore[T]) Save(ctx context.Context, id string, snap []byte) error {
	_, err := s.kv.Put(ctx, id, snap)
	return err
}

func (s *snapshotStore[T]) Load(ctx context.Context, id string) ([]byte, error) {
	v, err := s.kv.Get(ctx, id)
	if err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, store.ErrNoSnapshot
		}
		return nil, err
	}

	return v.Value(), nil
}
