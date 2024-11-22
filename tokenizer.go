package main

import "strings"

func tokenize(content []byte) []string {
	tokens := make([]string, 0)

	start := 0
	cursor := 0
	size := len(content)

outer:
	for {
		for isWhitespace(content[cursor]) {
			cursor++

			if cursor >= size {
				break outer
			}
		}

		if cursor >= size {
			break
		}

		start = cursor

		if isAlphaNumeric(content[cursor]) {
			for isAlphaNumeric(content[cursor]) {
				cursor++

				if cursor >= size {
					token := strings.ToLower(string(content[start:]))

					tokens = append(tokens, token)

					break outer
				}
			}

			token := strings.ToLower(string(content[start:cursor]))

			tokens = append(tokens, token)
			continue
		}

		c := content[cursor]

		if content[cursor] == c {
			for c == content[cursor] && cursor < size-1 {
				cursor++

				if cursor >= size {
					token := strings.ToLower(string(content[start:]))

					tokens = append(tokens, token)
					break outer
				}
			}

			token := strings.ToLower(string(content[start:cursor]))

			tokens = append(tokens, token)

		} else {
			token := strings.ToLower(string(content[start:cursor]))

			tokens = append(tokens, token)

			cursor++
		}
	}

	return tokens
}
