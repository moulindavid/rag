package chunker

import "strings"

const (
	maxChunkSize = 2000
	overlapSize  = 200
)

func Chunk(text string) []string {
	paragraphs := strings.Split(text, "\n\n")

	var result []string
	var buffer string

	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		if len(buffer) == 0 {
			buffer = p
		} else if len(buffer)+2+len(p) < maxChunkSize {
			buffer += "\n\n" + p
		} else {
			result = append(result, buffer)
			overlap := buffer
			if len(buffer) > overlapSize {
				overlap = buffer[len(buffer)-overlapSize:]
			}
			buffer = overlap + "\n\n" + p
		}
	}

	if buffer != "" {
		result = append(result, buffer)
	}

	return result
}
