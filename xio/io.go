package io

import (
	"io"
)

// ProduceAndConsume connect producer and consumer using pipe.
func ProduceAndConsume(producer func(io.Writer) error, consumer func(io.Reader) error) error {
	pr, pw := io.Pipe()
	go func() {
		if err := producer(pw); err != nil {
			pw.CloseWithError(err)
		} else {
			pw.Close()
		}
	}()

	return consumer(pr)
}
