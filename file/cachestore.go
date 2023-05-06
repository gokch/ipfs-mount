package file

import (
	"context"
	"errors"
	"time"

	theine "github.com/Yiling-J/theine-go"
	datastore "github.com/ipfs/go-datastore"
	query "github.com/ipfs/go-datastore/query"
)

type Cachestore struct {
	cache *theine.Cache[string, []byte]
	ttl   time.Duration
}

var _ datastore.Datastore = (*Cachestore)(nil)
var _ datastore.Batching = (*Cachestore)(nil)

func NewCacheStore(ttl time.Duration) *Cachestore {
	cache, _ := theine.NewBuilder[string, []byte](1000).Build()
	return &Cachestore{
		cache: cache,
		ttl:   ttl,
	}
}

func (ds *Cachestore) Put(ctx context.Context, key datastore.Key, value []byte) error {
	ds.cache.SetWithTTL(key.String(), value, 0, ds.ttl)
	return nil
}

func (ds *Cachestore) Sync(ctx context.Context, prefix datastore.Key) error {
	return nil
}

func (ds *Cachestore) Get(ctx context.Context, key datastore.Key) (value []byte, err error) {
	value, success := ds.cache.Get(key.String())
	if !success {
		return nil, datastore.ErrNotFound
	}
	return value, nil
}

func (ds *Cachestore) Has(ctx context.Context, key datastore.Key) (exists bool, err error) {
	val, _ := ds.cache.Get(key.String())
	return val != nil, nil
}

func (ds *Cachestore) GetSize(ctx context.Context, key datastore.Key) (size int, err error) {
	value, _ := ds.cache.Get(key.String())
	if value == nil {
		return -1, datastore.ErrNotFound
	}
	return len(value), nil
}

func (ds *Cachestore) Delete(ctx context.Context, key datastore.Key) (err error) {
	ds.cache.Delete(key.String())
	return nil
}

func (ds *Cachestore) Query(ctx context.Context, q query.Query) (query.Results, error) {
	return nil, errors.New("TODO implement query for rueidis datastore?")
}

func (ds *Cachestore) Batch(ctx context.Context) (datastore.Batch, error) {
	return nil, datastore.ErrBatchUnsupported
}

func (ds *Cachestore) Close() error {
	ds.cache.Close()
	return nil
}