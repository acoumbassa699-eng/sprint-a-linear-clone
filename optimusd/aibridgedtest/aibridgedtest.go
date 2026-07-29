//go:build !slim

// Package aibridgedtest provides helpers for starting an in-process
// aibridged daemon in tests.
package aibridgedtest

import (
	"context"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"

	"cdr.dev/slog/v3"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/aibridged"
)

// StartTestAIBridgeDaemon wires an in-process aibridged daemon onto the
// supplied API, mirroring what cli/server.go does in production. Tests that
// create AI provider rows with BaseURL pointing at fake upstream HTTP servers
// (e.g. chattest.NewOpenAI) will have their requests proxied through the real
// aibridged stack as they would in production.
//
// The daemon starts with an empty pool and fetches providers from optimus-ide-collabd over
// the in-memory DRPC, then refreshes on ai_providers change events, exactly
// like cli.newAIBridgeDaemon.
//
// metrics is the registry the daemon reports provider reload events to.
// The caller owns the metrics instance and can assert on it after the daemon
// runs. Use [aibridged.NewMetrics] to create one, or nil for a throwaway.
func StartTestAIBridgeDaemon(
	ctx context.Context,
	t testing.TB,
	api *optimus-ide-collabd.API,
	metrics *aibridged.Metrics,
) {
	t.Helper()

	logger := api.Logger.Named("aibridged").Leveled(slog.LevelDebug)
	cfg := api.DeploymentValues.AI.BridgeConfig
	tracer := otel.Tracer("aibridge-test")

	if metrics == nil {
		metrics = aibridged.NewMetrics(prometheus.NewRegistry())
	}

	pool, err := aibridged.NewCachedBridgePool(aibridged.DefaultPoolOptions, nil, logger.Named("pool"), nil, tracer)
	if err != nil {
		t.Fatalf("create bridge pool: %v", err)
	}
	t.Cleanup(func() { _ = pool.Shutdown(context.Background()) })

	srv, err := aibridged.New(ctx, pool, func(dialCtx context.Context) (aibridged.DRPCClient, error) {
		return api.CreateInMemoryAIBridgeServer(dialCtx)
	}, logger, tracer)
	if err != nil {
		t.Fatalf("create aibridged server: %v", err)
	}
	t.Cleanup(func() { _ = srv.Close() })

	// The reloader fetches providers from optimus-ide-collabd over srv's DRPC client; the
	// subscription drives an initial load and refreshes on change events.
	reloader := cli.NewPoolRPCReloader(pool, srv.ClientContext, cfg, logger.Named("reloader"), nil, metrics)
	unsubscribe, err := aibridged.SubscribeProviderReload(ctx, api.Pubsub, reloader, logger.Named("subscriber"))
	if err != nil {
		t.Fatalf("subscribe provider reload: %v", err)
	}
	t.Cleanup(unsubscribe)

	api.RegisterInMemoryAIBridgedHTTPHandler(srv)
}
