# Errors Channel

Implements `<-chan error` pattern and plumbing for error propagation

Import package with go mod, add to `go.mod`:
```
require github.com/artyomturkin/go-errors v1.0.0
```

Example:
```go
// create error channel
errChan := &errors.Channel{}

// create subscription
errCh := errChan.Errors()

// print errors
go func(){
    for err := range errCh {
        fmt.Printf("%v\n", err)
    }
}()

// publish an error
errChan.Publish(fmt.Errorf("new error"))

// close channel
errChan.Close()

// Output: new error
```