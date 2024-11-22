package tokenizer

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func chr(reader *bufio.Reader) rune {
	c, _, err := reader.ReadRune()

	if err == io.EOF {
		return '\000'
	} else if err != nil {
		panic(err)
	}

	return c
}

func Tokenize(file *os.File) []string {
	tokens := make([]string, 0)

	reader := bufio.NewReader(file)

	c := chr(reader)

outer:
	for {
		for isWhitespace(c) {
			c = chr(reader)

			if c == '\000' {
				break outer
			}
		}

		if c == '\000' {
			break
		}

		chars := make([]rune, 0)

		if isAlphaNumeric(c) {
			for isAlphaNumeric(c) {
				chars = append(chars, c)

				c = chr(reader)

				if c == '\000' {
					token := strings.ToLower(string(chars))

					tokens = append(tokens, token)

					break outer
				}
			}

			token := strings.ToLower(string(chars))

			tokens = append(tokens, token)
			continue
		}

		tmp := c

		if c == tmp {
			for c == tmp && c != '\000' {
				chars = append(chars, c)

				c = chr(reader)

				if c == '\000' {
					token := strings.ToLower(string(chars))

					tokens = append(tokens, token)
					break outer
				}
			}

			token := strings.ToLower(string(chars))

			tokens = append(tokens, token)
		} else {
			chars = append(chars, c)

			token := strings.ToLower(string(chars))

			tokens = append(tokens, token)

			c = chr(reader)
		}
	}

	return tokens
}
