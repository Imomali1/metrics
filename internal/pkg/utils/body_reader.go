package utils

import "io"

func ReadAll(r io.Reader) ([]byte, error) {
	// Максимальное количество байтов нужно изменить,
	// если объём данных будет большим
	const maxBodySize = 10 * 1024
	b := make([]byte, 0, maxBodySize)
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
