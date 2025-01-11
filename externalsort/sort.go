//go:build !change

package externalsort

import (
	"bufio"
	"io"
)

func NewReader(r io.Reader) LineReader {
	return &ReaderImpl{
		ioReader:    r,
		bufioReader: bufio.NewReader(r),
		lineEnding:  LineEnding.Unix,
	}
}

func NewWriter(w io.Writer) LineWriter {
	return &WriterImpl{
		ioWriter:   w,
		lineEnding: LineEnding.Unix,
	}
}

func Merge(w LineWriter, readers ...LineReader) error {
	h, err := initHeap(readers)
	if err != nil {
		return err
	}

	return mergeLines(w, h)
}

func Sort(w io.Writer, in ...string) error {
	var err error
	var lineEnding LineEndingType

	for _, filename := range in {
		err = processSingleFile(filename, &lineEnding)
		if err != nil {
			return err
		}
	}

	return mergeFiles(w, lineEnding, in...)
}
