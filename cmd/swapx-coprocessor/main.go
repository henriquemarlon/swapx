package main

import (
	"context"
	"log/slog"

	"github.com/henriquemarlon/swapx/pkg/router"
	"github.com/rollmelette/rollmelette"
)

func Echo(
	env rollmelette.Env,
	metadata rollmelette.Metadata,
	deposit rollmelette.Deposit,
	payload []byte,
) error {
	env.Voucher(metadata.MsgSender, payload)
	env.Notice(payload)
	env.Report(payload)
	return nil
}

func NewDApp(echo router.AdvanceHandlerFunc) *router.Router {
	r := router.NewRouter()
	r.HandleAdvance("echo", echo)
	return r
}

func main() {
	//////////////////////// Setup DApp /////////////////////////
	app := NewDApp(Echo)
	ctx := context.Background()
	opts := rollmelette.NewRunOpts()
	if err := rollmelette.Run(ctx, opts, app); err != nil {
		slog.Error("application error", "error", err)
	}
}
