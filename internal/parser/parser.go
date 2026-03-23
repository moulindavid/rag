package parser

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/dslipak/pdf"
)

func Parse(filename string, reader io.Reader) (string, error) {
	extension := filepath.Ext(filename)
	switch extension {
	case ".txt":
		file, err := io.ReadAll(reader)
		if err != nil {
			return "", err
		}
		return normalizeLineEndings(strings.TrimSpace(string(file))), nil

	case ".pdf":
		data, err := io.ReadAll(reader)
		if err != nil {
			return "", fmt.Errorf("reading pdf: %w", err)
		}

		r, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return "", fmt.Errorf("parsing pdf: %w", err)
		}

		var sb strings.Builder
		for i := 1; i <= r.NumPage(); i++ {
			page := r.Page(i)
			text, err := page.GetPlainText(nil)
			if err != nil {
				return "", fmt.Errorf("extracting text from page %d: %w", i, err)
			}
			sb.WriteString(text)
		}
		return normalizeLineEndings(strings.TrimSpace(sb.String())), nil

	default:
		return "", fmt.Errorf("unsupported file extension: %s", extension)
	}
}

func normalizeLineEndings(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}
