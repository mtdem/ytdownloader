package video

import (
	"io"

	"github.com/gosuri/uiprogress"
)

// ProgressBarReader wraps an existing io.Reader.
// It simply forwards the Read() call, while tracking
// progress to progress bar.
type ProgressBarReader struct {
	io.Reader
	total int64 // Total of bytes transferred
	bar   *uiprogress.Bar
}

// Read 'overrides' the underlying io.Reader's Read method.
// This is the one that will be called by io.Copy().
// We use it to keep track progress using bar and then forward the call.
func (pbr *ProgressBarReader) Read(p []byte) (int, error) {
	n, err := pbr.Reader.Read(p)
	pbr.total += int64(n)
	pbr.bar.Set(int(pbr.total) + 1)
	return n, err
}
