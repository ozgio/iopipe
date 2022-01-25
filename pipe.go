package iopipe

import (
	"context"
	"io"
	"sync"
)

type Writer func(context.Context, io.Writer) error
type Reader func(context.Context, io.Reader) error

// Parallel creates a synchronous in memory pipe, runs reader and writer functions as go routines
// and waits until both of them are quits. It closes the reader and writer at the end of functions.
// ctx is cancelled automatically if any of the functions  returns an error
//
// It uses io.Pipe under the hood, see https://pkg.go.dev/io#Pipe for details.
func Parallel(ctx context.Context, write Writer, read Reader) error {
	pr, pw := io.Pipe()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	var once sync.Once
	var firstErr error
	do := func(fn func() error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := fn()
			if err != nil {
				once.Do(func() {
					firstErr = err
					cancel()
				})
			}
		}()
	}

	do(func() error {
		defer pw.Close()
		return write(ctx, pw)
	})

	do(func() error {
		defer pr.Close()
		return read(ctx, pr)
	})

	wg.Wait()
	return firstErr
}
