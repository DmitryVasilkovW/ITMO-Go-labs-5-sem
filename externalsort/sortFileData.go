package externalsort

import (
	"errors"
	"io"
	"os"
	"slices"
	"strings"
)

const accessRights = 0600

func processSingleFile(filename string, lineEnding *LineEndingType) error {
	lines, err := readFileLines(filename, lineEnding)
	if err != nil {
		return err
	}

	if len(lines) > 1 {
		slices.SortFunc(lines, strings.Compare)
	}

	return writeFileLines(filename, *lineEnding, lines)
}

func readFileLines(filename string, lineEnding *LineEndingType) ([]string, error) {
	file, err := openFile(filename)
	if err != nil {
		return nil, err
	}

	defer func(f *os.File) {
		err = f.Close()
	}(file)

	err = detectAndSetLineEnding(file, lineEnding)
	if err != nil {
		return nil, err
	}

	return readLines(file, *lineEnding)
}

func openFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, accessRights)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func detectAndSetLineEnding(f *os.File, lineEnding *LineEndingType) error {
	if len(*lineEnding) == 0 {
		detectedLineEnding, err := detectLineEnding(f, -1)
		if err != nil {
			return err
		}
		*lineEnding = detectedLineEnding
	}
	return nil
}

func readLines(f *os.File, lineEnding LineEndingType) ([]string, error) {
	var lines []string
	reader := NewReader(f).(*ReaderImpl).changeLineEnding(lineEnding)

	for {
		line, err := reader.ReadLine()
		if err == nil || (errors.Is(err, io.EOF) && len(line) > 0) {
			lines = append(lines, line)
		}
		if err != nil {
			break
		}
	}
	return lines, nil
}

func writeFileLines(filename string, lineEnding LineEndingType, lines []string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY, accessRights)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err = f.Close()
	}(file)

	writer := NewWriter(file).(*WriterImpl).changeLineEnding(lineEnding)
	for _, line := range lines {
		err = writer.Write(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func mergeFiles(w io.Writer, lineEnding LineEndingType, in ...string) error {
	readers := make([]LineReader, 0, len(in))
	writer := NewWriter(w).(*WriterImpl).changeLineEnding(lineEnding)

	for _, filename := range in {
		file, err := os.OpenFile(filename, os.O_RDONLY, accessRights)
		if err != nil {
			return err
		}
		//goland:noinspection GoDeferInLoop
		defer func(f *os.File) {
			err = f.Close()
		}(file)

		reader := NewReader(file).(*ReaderImpl).changeLineEnding(lineEnding)
		readers = append(readers, reader)
	}

	return Merge(writer, readers...)
}
