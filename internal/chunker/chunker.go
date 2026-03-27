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
				start := len(buffer) - overlapSize
				for start > 0 && buffer[start-1] != ' ' {
					start--
				}
				overlap = buffer[start:]
			}
			buffer = overlap + "\n\n" + p
		}
	}

	if buffer != "" {
		result = append(result, buffer)
	}

	return result
}
