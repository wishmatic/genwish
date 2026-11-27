package cli

import (
	"context"
	"encoding/base64"

	"github.com/wishmatic/genwish/internal/resolve"
)

func DataURL(image resolve.Image) string {
	return "data:" + image.MediaType + ";base64," + base64.StdEncoding.EncodeToString(image.Data)
}

func ResolveDataURL(ctx context.Context, resolver *resolve.Client, address string) (string, error) {
	if address == "" {
		return "", nil
	}

	image, err := resolver.Resolve(ctx, address)
	if err != nil {
		return "", err
	}

	return DataURL(image), nil
}

func ResolveDataURLs(
	ctx context.Context, resolver *resolve.Client, addresses []string,
) ([]string, error) {
	if len(addresses) == 0 {
		return nil, nil
	}

	dataURLs := make([]string, 0, len(addresses))

	for _, address := range addresses {
		dataURL, err := ResolveDataURL(ctx, resolver, address)
		if err != nil {
			return nil, err
		}

		dataURLs = append(dataURLs, dataURL)
	}

	return dataURLs, nil
}
