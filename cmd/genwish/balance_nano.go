package main

import (
	"context"
	"fmt"

	"github.com/wishmatic/genwish/internal/cli"
)

func runBalanceNano(ctx context.Context, deps cli.Deps) error {
	client, err := cli.NanoClient(deps)
	if err != nil {
		return err
	}

	balance, err := client.CheckBalance(ctx)
	if err != nil {
		return fmt.Errorf("balance nano: %w", err)
	}

	fmt.Fprintf(deps.Out, "USD: %s\n", cli.FormatUSD(balance.USD))
	fmt.Fprintf(deps.Out, "Nano: %s\n", cli.FormatAmount(balance.Nano))

	return nil
}
