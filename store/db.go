package store

// TODO: move out of rpc, since you may have other modules (e.g. scheduled tasks) that want to use the database
// See example in hvm_listings

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cameront/template_backend/config"
	"github.com/cameront/template_backend/ent"
	"github.com/cameront/template_backend/logging"
)

/*

It might be nice for some applications to separate read/write connections.
Especially in situations where you want consistend reads (e.g. don't want to see
new rows in between queries of the same table that were added by a committed
transaction by some other thread).

In WAL mode, you can have deferred transactions that just keep your reads
consistent. You can only have one writer (which should use BEGIN IMMEDIATE to
ensure that only one writer is attempting to acquire the lock at any one time)

So you'd basically initialize twice, once for reads and once for writes

read_params:  "?mode=ro&_txlock=deferred"
write_params: "?mode=rw&_journal=WAL&_timeout=5000&_fk=true&_sync=NORMAL&_txlock=immediate"

And you'd have your DbProvider make two functions:
ReadTx(ctx, ...) which uses the read client
WriteTx(ctx, ...) which uses the write client

These would replace the WithTransaction function today

See also:
 https://github.com/mattn/go-sqlite3/issues/1022#issuecomment-1067353980
 https://phiresky.github.io/blog/2020/sqlite-performance-tuning/
 https://github.com/mattn/go-sqlite3/issues/1179

*/

func InitStore(ctx context.Context) (*clientWrapper, error) {
	cfg := config.MustContext(ctx)

	entClient, err := ent.Open(cfg.DB_DriverName, cfg.DB_URI)
	if err != nil {
		return nil, err
	}

	return &clientWrapper{dbClient: entClient}, nil
}

type ctxKey string

const (
	txKey   ctxKey = "tx"
	inTxKey ctxKey = "inTx"
)

func WithTransaction[O any](ctx context.Context, db DbProvider, fn func(context.Context, *ent.Client) (out O, err error)) (O, error) {
	dummyO := new(O)

	client := db.Get(ctx)
	if inTx, ok := ctx.Value(inTxKey).(bool); ok && inTx {
		// We're already inside of a transaction, so just run the function.
		return fn(ctx, client)
	}

	tx, err := client.BeginTx(ctx, nil)
	if err != nil {
		return *dummyO, err
	}

	txClient := tx.Client()
	ctx = context.WithValue(ctx, txKey, txClient)
	ctx = context.WithValue(ctx, inTxKey, true)

	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()

	out, err := fn(ctx, txClient)
	if err != nil {
		logging.Errorf(ctx, "rolling back transaction due to error: %v", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			// hmm, likely sorts screwed here!?
			slog.Error(fmt.Sprintf("rollback err: %v", rollbackErr))
		}
		return *dummyO, err
	}
	err = tx.Commit()
	return out, err
}

// We should always retrieve the db client from context, so that we don't leak
// outside our own transactions. This goofy little interface helps enforce that.
type DbProvider interface {
	// Get() is used by callers that don't care whether they're inside of a txn
	Get(ctx context.Context) (client *ent.Client)
	// GetTx() is used by callers that do care whether they're inside of a txn
	GetTx(ctx context.Context, txnRequiredReason string) (client *ent.Client, err error)
}

type clientWrapper struct {
	dbClient *ent.Client
}

func (c *clientWrapper) Close() {
	c.dbClient.Close()
}

func (c *clientWrapper) Get(ctx context.Context) *ent.Client {
	client, _ := c.GetTx(ctx, "")
	return client
}

func (c *clientWrapper) GetTx(ctx context.Context, txnRequiredReason string) (client *ent.Client, err error) {
	// First check for transactional client (set by WithTransaction)
	txI := ctx.Value("tx")
	if txI != nil {
		tx, ok := txI.(*ent.Client)
		if !ok {
			return nil, fmt.Errorf("context transaction exists, but expected type *ent.Client, got: %T", tx)
		}
		return tx, nil
	}
	if txnRequiredReason != "" {
		return nil, fmt.Errorf("no txn found, but required for: %s", txnRequiredReason)
	}
	return c.dbClient, nil
}
