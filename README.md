# Parallel IO Pipe

iopipe provides the function Parallel to create a synchronous in memory pipe and lets you write to 
and read from the pipe parallely. 

## Usage

```go
// Writes `Hello world` to standard output
err := iopipe.Parallel(context.Background(), func(ctx context.Context, w io.Writer) error {
    _, err := fmt.Fprint(w, "Hello World")
    return err
}, func(ctx context.Context, r io.Reader) error {
    _, err := io.Copy(os.Stdout, pr)
    return err
})
```