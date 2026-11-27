package balance

import (
	"context"
	"fmt"

	"github.com/wishmatic/genwish/internal/cli"
)

func RunNovel(ctx context.Context, deps cli.Deps) error {
	client, err := cli.NovelClient(deps)
	if err != nil {
		return err
	}

	balance, err := client.CheckBalance(ctx)
	if err != nil {
		return fmt.Errorf("balance novel: %w", err)
	}

	fmt.Fprintf(deps.Out, "Anlas: %d\n", balance.Anlas)

	if balance.Tier != "" {
		fmt.Fprintf(deps.Out, "Tier: %s\n", balance.Tier)
	}

	return nil
}
