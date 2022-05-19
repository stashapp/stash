package transport

import (
	"io"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytesRead(t *testing.T) {
	t.Run("test concurrency", func(t *testing.T) {
		// Test for the race detector, to verify a Read that doesn't yield any bytes
		// is okay to use from multiple goroutines. This was our historic behavior.
		// See golang.org/issue/7856
		r := bytesReader{s: &([]byte{})}
		var wg sync.WaitGroup
		for i := 0; i < 5; i++ {
			wg.Add(2)
			go func() {
				defer wg.Done()
				var buf [1]byte
				r.Read(buf[:])
			}()
			go func() {
				defer wg.Done()
				r.Read(nil)
			}()
		}
		wg.Wait()
	})

	t.Run("fail to read if pointer is nil", func(t *testing.T) {
		n, err := (&bytesReader{}).Read(nil)
		require.Equal(t, 0, n)
		require.NotNil(t, err)
		require.Equal(t, "byte slice pointer is nil", err.Error())
	})

	t.Run("read using buffer", func(t *testing.T) {
		data := []byte("0123456789")
		r := bytesReader{s: &data}

		got := make([]byte, 0, 11)
		buf := make([]byte, 1)
		for {
			n, err := r.Read(buf)
			if n < 0 {
				require.Fail(t, "unexpected bytes read size")
			}
			got = append(got, buf[:n]...)
			if err != nil {
				if err == io.EOF {
					break
				}
				require.Fail(t, "unexpected error while reading", err.Error())
			}
		}
		require.Equal(t, "0123456789", string(got))
	})

	t.Run("read updated pointer value", func(t *testing.T) {
		data := []byte("0123456789")
		pointer := &data
		r := bytesReader{s: pointer}
		data[2] = []byte("9")[0]

		got := make([]byte, 0, 11)
		buf := make([]byte, 1)
		for {
			n, err := r.Read(buf)
			if n < 0 {
				require.Fail(t, "unexpected bytes read size")
			}
			got = append(got, buf[:n]...)
			if err != nil {
				if err == io.EOF {
					break
				}
				require.Fail(t, "unexpected error while reading", err.Error())
			}
		}
		require.Equal(t, "0193456789", string(got))
	})
}
