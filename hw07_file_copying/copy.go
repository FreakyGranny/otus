package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	// ErrUnsupportedFile ...
	ErrUnsupportedFile = errors.New("unsupported file")
	// ErrOffsetExceedsFileSize ...
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

// Copy ...
func Copy(fromPath string, toPath string, offset, limit int64) error {
	info, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if info.IsDir() || !info.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	if offset >= info.Size() {
		return ErrOffsetExceedsFileSize
	}

	src, err := os.OpenFile(fromPath, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	dst, err := os.OpenFile(toPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	bytesLeft := info.Size() - offset
	barSize := limit
	if bytesLeft < limit || limit == 0 {
		barSize = bytesLeft
	}

	bar := pb.Full.Start64(barSize)
	defer bar.Finish()
	barReader := bar.NewProxyReader(src)

	_, err = src.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	if limit == 0 {
		_, err = io.Copy(dst, barReader)
	} else {
		_, err = io.CopyN(dst, barReader, limit)
	}
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}
