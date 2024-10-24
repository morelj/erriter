# Erriter

Package erriter provide a way to create iterators which can generate errors during iteration.

It is compatible with Go 1.23 iter package.

This package defines `Seq` and `Seq2` types which are wrappers around standard iterators which are able to store an
error which may occur internally during iteration.

Iterators are created using the `New` and `New2` functions which both take an iterator function
(similar to iterator functions defined by the standard `iter` package), which return an error.

The `Seq.Range` and `Seq2.Range` methods can then be used as iterator functions.

Errors can be checked after the iterator has returned by calling `Seq.Error` or `Seq2.Error`.

## Example

```go
// Create an erriter.Seq
iter := erriter.New(func (yield func(int) error) {
    if !yield(1) {
        return false
    }
    return errors.New("something went wrong...")
})

// Iterate
for v := range iter.Range {
    fmt.Println(v)
}

// Check for errors
if err := iter.Error(); err != nil {
    panic(err)
}
```

## Installation

Install using go get:

```bash
go get github.com/morelj/erriter
```
