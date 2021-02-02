package main

import (
	"context"

	"github.com/fguy/dispatch-simulation/config"
	"github.com/fguy/dispatch-simulation/services/courier"
	"github.com/fguy/dispatch-simulation/services/dispatch"
	"github.com/fguy/dispatch-simulation/services/fulfillment"
	"github.com/fguy/dispatch-simulation/utils"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// RegisterFulfillment starts the simulation.
//
// Register is a typical top-level application function: it takes a generic
// type like Server, which typically comes from a third-party library, and
// introduces it to a type that contains our application logic. In this case,
// that introduction consists of registering a fulfillment service. Other typical
// examples include registering RPC procedures and starting queue consumers.
//
// Fx calls these functions invocations, and they're treated differently from
// the constructor functions above. Their arguments are still supplied via
// dependency injection and they may still return an error to indicate
// failure, but any other return values are ignored.
//
// Unlike constructors, invocations are called eagerly. See the main function
// below for details.
func RegisterFulfillment(
	lc fx.Lifecycle,
	logger *zap.Logger,
	fulfillmentService fulfillment.Interface,
) {
	// If fulfillmentService.Start() is called, we know that another function is using the process. In
	// that case, we'll use the Lifecycle type to register a Hook that starts
	// and stops our simulation.
	//
	// Hooks are executed in dependency order. At startup, NewLogger's hooks run
	// before fulfillmentService.Start. On shutdown, the order is reversed.
	//
	// Returning an error from OnStart hooks interrupts application startup. Fx
	// immediately runs the OnStop portions of any successfully-executed OnStart
	// hooks (so that types which started cleanly can also shut down cleanly),
	// then exits.
	//
	// Returning an error from OnStop hooks logs a warning, but Fx continues to
	// run the remaining hooks.
	lc.Append(fx.Hook{
		// To mitigate the impact of deadlocks in application startup and
		// shutdown, Fx imposes a time limit on OnStart and OnStop hooks. By
		// default, hooks have a total of 30 seconds to complete. Timeouts are
		// passed via Go's usual context.Context.
		OnStart: func(context.Context) error {
			logger.Info("Starting the simulation")
			return fulfillmentService.Start()
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping the simulation.")
			return fulfillmentService.Stop()
		},
	})
}

var app = fx.New(
	// Provide all the constructors we need, which teaches Fx how we'd like to
	// construct the *config.AppConfig, *zap.Logger and *http.Server types.
	// Remember that constructors are called lazily, so this block doesn't do
	// much on its own.
	fx.Provide(
		config.NewAppConfig,
		courier.New,
		dispatch.New,
		fulfillment.New,
		NewLogger,
		utils.NewCourierQueue,
		utils.NewOrderQueue,
		utils.NewStat,
	),
	// Since constructors are called lazily, we need some invocations to
	// kick-start our application. In this case, we'll use Register. Since it
	// depends on an fulfillment.Interface, calling it requires Fx
	// to build the type using the constructors above. Since we call
	// fulfillment.New, we also register Lifecycle hooks to start and stop
	fx.Invoke(
		RegisterFulfillment,
	),
)

func main() {
	app.Run()
}
