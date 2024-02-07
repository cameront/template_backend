package counterservice

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cameront/go-svelte-sqlite-template/config"
	"github.com/cameront/go-svelte-sqlite-template/ent"
	"github.com/cameront/go-svelte-sqlite-template/logging"
)

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

func WithTransaction[O any](ctx context.Context, db dbProvider, fn func(context.Context, *ent.Client) (out O, err error)) (O, error) {
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
		logging.GetLogger(ctx).Error("rolling back transaction")
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
type dbProvider interface {
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
