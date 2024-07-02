package utils

import "io"

// MaxBodySize - max count of bytes that reader can read.
const MaxBodySize = 10 * 1024

// ReadAll reads bytes from reader.
func ReadAll(r io.Reader) ([]byte, error) {
	b := make([]byte, 0, MaxBodySize)
	for {
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b[:len(b):len(b)], err
		}

		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
	}
}
