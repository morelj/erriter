// Package erriter provide a way to create iterators which can generate errors during iteration.
//
// It is compatible with Go 1.23 iter package.
//
// This package defines [Seq] and [Seq2] types which are wrappers around standard iterators which
// are able to store an error which may occur internally during iteration.
//
// Iterators are created using the [New] and [New2] functions which both take an iterator function
// (similar to iterator functions defined by the standard iter package), which return an error.
//
// The [Seq.Range] and [Seq2.Range] methods can then be used as iterator functions.
//
// Errors can be checked after the iterator has returned by calling [Seq.Error] or [Seq2.Error].
//
//	// Create an erriter.Seq
//	iter := erriter.New(func (yield func(int) error) {
//	  if !yield(1) {
//	    return false
//	  }
//	  return errors.New("something went wrong...")
//	})
//	// Iterate
//	for v := range iter.Range {
//	  fmt.Println(v)
//	}
//	// Check for errors
//	if err := iter.Error(); err != nil {
//	  panic(err)
//	}
//
// Each time the [Seq.Range] or [Seq2.Range] method is called, the internal error is reset.
package erriter

// A Seq wraps an iterator function which yields one value and an internal error.
//
// The [Seq.Range] is an iterator function which calls the [Seq]'s iterator function and updates its internal error
// to the return value of it.
type Seq[V any] struct {
	f   func(yield func(V) bool) error
	err error
}

// New creates a new [Seq] with the given iterator function.
//
// The returned [Seq]'s internal error is set to nil.
func New[V any](f func(yield func(V) bool) error) *Seq[V] {
	return &Seq[V]{f: f}
}

// Range is an iterator function which will call the iterator function of s and update its error according to
// the iterator's return value.
func (s *Seq[V]) Range(yield func(V) bool) {
	s.err = s.f(yield)
}

// Error returns the error of s. This should be called after [Seq.Range] has been called.
func (s *Seq[V]) Error() error {
	return s.err
}

// A Seq2 wraps an iterator function which yields two values and an internal error.
//
// The [Seq2.Range] is an iterator function which calls the [Seq2]'s iterator function and updates its internal error
// to the return value of it.
type Seq2[K, V any] struct {
	f   func(yield func(K, V) bool) error
	err error
}

// New2 creates a new [Seq2] with the given iterator function.
//
// The returned [Seq2]'s internal error is set to nil.
func New2[K, V any](f func(yield func(K, V) bool) error) *Seq2[K, V] {
	return &Seq2[K, V]{f: f}
}

// Range is an iterator function which will call the iterator function of s and update its error according to
// the iterator's return value.
func (s *Seq2[K, V]) Range(yield func(K, V) bool) {
	s.err = s.f(yield)
}

// Error returns the error of s. This should be called after [Seq2.Range] has been called.
func (s *Seq2[K, V]) Error() error {
	return s.err
}
