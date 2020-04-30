package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

const (
	bufsize = 1 << 10
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	f, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	t, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	defer t.Close()
	fs, _ := f.Stat()
	if fs.Size() < offset {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 || limit+offset > fs.Size() {
		limit = fs.Size() - offset
	}
	o, err := f.Seek(offset, 0)
	if err != nil || o != offset {
		return err
	}
	fmt.Printf("Coping file %v to file %v", fromPath, toPath)
	fmt.Println()
	bar := pb.StartNew(int(limit))
	for i := offset; i < offset+limit; i += bufsize {
		copylen := int64(math.Min(bufsize, float64(limit)))
		c, err := io.CopyN(t, f, copylen)
		if (err != nil && err != io.EOF) || c > limit {
			return err
		}
		printProgressBar(bar, int(c))
	}
	bar.Finish()
	return nil
}

func printProgressBar(bar *pb.ProgressBar, lastchunk int) {
	bar.Add(lastchunk)
	time.Sleep(time.Millisecond) //для визуализации Progressbar
}
