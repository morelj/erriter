package erriter

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeq(t *testing.T) {
	type result struct {
		expected []int
		error    bool
	}

	cases := []struct {
		seq      *Seq[int]
		runCount int
		results  []result
	}{
		{
			seq: New(func(yield func(int) bool) error {
				for _, i := range []int{1, 2, 3} {
					if !yield(i) {
						return nil
					}
				}
				return nil
			}),
			runCount: 1,
			results: []result{
				{
					expected: []int{1, 2, 3},
				},
			},
		},
		{
			seq: New(func(yield func(int) bool) error {
				return errors.New("error")
			}),
			runCount: 1,
			results: []result{
				{
					error: true,
				},
			},
		},
		{
			seq: func() *Seq[int] {
				run := -1
				return New(func(yield func(int) bool) error {
					run++
					switch run {
					case 0:
						// First run, return an error
						return errors.New("error")
					case 1:
						// Second run, return something
						yield(42)
						return nil
					default:
						panic("Too many runs!")
					}
				})
			}(),
			results: []result{
				{
					error: true,
				},
				{
					expected: []int{42},
				},
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			for run := 0; run < c.runCount; run++ {
				res := make([]int, 0)
				for value := range c.seq.Range {
					res = append(res, value)
				}
				if c.results[run].error {
					assert.Errorf(c.seq.Error(), "Run #%d", run)
				} else {
					require.NoErrorf(c.seq.Error(), "Run #%d", run)
					assert.Equalf(c.results[run].expected, res, "Run #%d", run)
				}
			}
		})
	}
}

func TestSeq2(t *testing.T) {
	type result struct {
		expected map[string]int
		error    bool
	}

	cases := []struct {
		seq      *Seq2[string, int]
		runCount int
		results  []result
	}{
		{
			seq: New2(func(yield func(string, int) bool) error {
				m := map[string]int{
					"one":   1,
					"two":   2,
					"three": 3,
				}
				for k, v := range m {
					if !yield(k, v) {
						return nil
					}
				}
				return nil
			}),
			runCount: 1,
			results: []result{
				{
					expected: map[string]int{
						"one":   1,
						"two":   2,
						"three": 3,
					},
				},
			},
		},
		{
			seq: New2(func(yield func(string, int) bool) error {
				return errors.New("error")
			}),
			runCount: 1,
			results: []result{
				{
					error: true,
				},
			},
		},
		{
			seq: func() *Seq2[string, int] {
				run := -1
				return New2(func(yield func(string, int) bool) error {
					run++
					switch run {
					case 0:
						// First run, return an error
						return errors.New("error")
					case 1:
						// Second run, return something
						yield("the answer", 42)
						return nil
					default:
						panic("Too many runs!")
					}
				})
			}(),
			results: []result{
				{
					error: true,
				},
				{
					expected: map[string]int{
						"the answer": 42,
					},
				},
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			for run := 0; run < c.runCount; run++ {
				res := map[string]int{}
				for k, v := range c.seq.Range {
					res[k] = v
				}
				if c.results[run].error {
					assert.Errorf(c.seq.Error(), "Run #%d", run)
				} else {
					require.NoErrorf(c.seq.Error(), "Run #%d", run)
					assert.Equalf(c.results[run].expected, res, "Run #%d", run)
				}
			}
		})
	}
}

func ExampleSeq() {
	iter := New(func(yield func(int) bool) error {
		for i, v := range []int{1, 2, -1, 3} {
			if v < 0 {
				return fmt.Errorf("negative value at index %d: %d", i, v)
			}
			if !yield(v) {
				return nil
			}
		}
		return nil
	})

	for v := range iter.Range {
		fmt.Printf("Got value: %d\n", v)
	}
	if err := iter.Error(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Output:
	// Got value: 1
	// Got value: 2
	// Error: negative value at index 2: -1
}
