package cli

import (
	"context"
	"encoding/base64"

	"github.com/wishmatic/genwish/internal/resolve"
)

func DataURL(image resolve.Image) string {
	return "data:" + image.MediaType + ";base64," + base64.StdEncoding.EncodeToString(image.Data)
}

// ImageDataURL turns an input image into a data URL. The input may be a local path, an http(s) URL,
// a data URI, or raw base64 data.
func ImageDataURL(ctx context.Context, resolver *resolve.Client, address string) (string, error) {
	if address == "" {
		return "", nil
	}

	image, err := readImageInput(ctx, resolver, address)
	if err != nil {
		return "", err
	}

	return DataURL(image), nil
}

func ImageDataURLs(
	ctx context.Context, resolver *resolve.Client, addresses []string,
) ([]string, error) {
	if len(addresses) == 0 {
		return nil, nil
	}

	dataURLs := make([]string, 0, len(addresses))

	for _, address := range addresses {
		dataURL, err := ImageDataURL(ctx, resolver, address)
		if err != nil {
			return nil, err
		}

		dataURLs = append(dataURLs, dataURL)
	}

	return dataURLs, nil
}

func readImageInput(
	ctx context.Context, resolver *resolve.Client, address string,
) (resolve.Image, error) {
	data, isFile, err := readLocalImage(address)
	if err != nil {
		return resolve.Image{}, err
	}

	if isFile {
		mediaType, err := resolve.Identify(data)
		if err != nil {
			return resolve.Image{}, err
		}

		return resolve.Image{Data: data, MediaType: mediaType}, nil
	}

	return resolver.Resolve(ctx, address)
}
