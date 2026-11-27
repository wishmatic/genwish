package cli

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wishmatic/genwish/internal/config"
	"github.com/wishmatic/genwish/internal/format"
)

func OutputFormat(cfg config.Config, requested string) (format.Format, error) {
	if requested != "" {
		return format.Parse(requested)
	}

	if cfg.Output.Format == "" {
		return format.Default, nil
	}

	return format.Parse(cfg.Output.Format)
}

func OutputDir(cfg config.Config, requested string) string {
	if strings.TrimSpace(requested) != "" {
		return requested
	}

	if cfg.Output.Path != "" {
		return cfg.Output.Path
	}

	return "."
}

func WriteImage(base, kind string, data []byte, mediaType string) (string, error) {
	ext, err := imageExtension(mediaType)
	if err != nil {
		return "", err
	}

	return writeMediaFile(base, kind, ext, data)
}

func WriteImages(base, kind string, images [][]byte, target format.Format) ([]string, error) {
	paths := make([]string, 0, len(images))

	for index, data := range images {
		converted, err := format.Convert(data, target)
		if err != nil {
			return nil, fmt.Errorf("image %d: %w", index+1, err)
		}

		path, err := WriteImage(base, kind, converted, target.MediaType())
		if err != nil {
			return nil, err
		}

		paths = append(paths, path)
	}

	return paths, nil
}

func PrintPaths(w io.Writer, paths []string) error {
	for _, path := range paths {
		if _, err := fmt.Fprintln(w, path); err != nil {
			return fmt.Errorf("write path: %w", err)
		}
	}

	return nil
}

func writeMediaFile(base, kind, ext string, data []byte) (string, error) {
	id, err := newID()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(base, kind, time.Now().UTC().Format("2006-01-02"))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create output directory: %w", err)
	}

	tmp, err := os.CreateTemp(dir, "."+id+"-*.tmp")
	if err != nil {
		return "", fmt.Errorf("create output file: %w", err)
	}

	name := tmp.Name()
	defer func() {
		_ = os.Remove(name)
	}()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()

		return "", fmt.Errorf("write output file: %w", err)
	}

	if err := tmp.Chmod(0o644); err != nil {
		_ = tmp.Close()

		return "", fmt.Errorf("set output permissions: %w", err)
	}

	if err := tmp.Close(); err != nil {
		return "", fmt.Errorf("close output file: %w", err)
	}

	final := filepath.Join(dir, id+"."+ext)

	if err := os.Rename(name, final); err != nil {
		return "", fmt.Errorf("store output file: %w", err)
	}

	name = ""

	return final, nil
}

func newID() (string, error) {
	var b [16]byte

	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("generate file name: %w", err)
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

func imageExtension(mediaType string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(mediaType))

	for _, candidate := range []format.Format{format.PNG, format.JPEG, format.JXL, format.WebP} {
		if candidate.MediaType() == normalized {
			return candidate.Extension(), nil
		}
	}

	return "", fmt.Errorf("cannot store media type %q", mediaType)
}
