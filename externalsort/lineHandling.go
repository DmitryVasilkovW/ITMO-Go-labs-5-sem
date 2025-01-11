package externalsort

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type LineEndingType []byte

var LineEnding = struct {
	Windows LineEndingType
	Unix    LineEndingType
}{
	Windows: LineEndingType{'\r', '\n'},
	Unix:    LineEndingType{'\n'},
}

func detectLineEnding(file *os.File, maxLinesToCheck int) (LineEndingType, error) {
	err := resetFilePosition(file)
	if err != nil {
		return nil, err
	}

	lineEndingCounts := map[string]uint64{"\r\n": 0, "\n": 0}
	bufferedReader := bufio.NewReader(file)

	fillMap(lineEndingCounts, maxLinesToCheck, bufferedReader)

	err = resetFilePosition(file)
	if err != nil {
		return nil, err
	}

	return determineMostFrequentLineEnding(lineEndingCounts), nil
}

func fillMap(hashMap map[string]uint64, maxLinesToCheck int, bufferedReader *bufio.Reader) {
	for lineIndex := 0; maxLinesToCheck < 0 || lineIndex < maxLinesToCheck; lineIndex++ {
		line, readErr := bufferedReader.ReadString('\n')

		if len(line) > 0 {
			switch {
			case strings.HasSuffix(line, "\r\n"):
				hashMap["\r\n"]++
			case strings.HasSuffix(line, "\n"):
				hashMap["\n"]++
			}
		}

		if readErr != nil {
			break
		}
	}
}

func resetFilePosition(file *os.File) error {
	_, err := file.Seek(0, io.SeekStart)
	return err
}

func determineMostFrequentLineEnding(lineEndingCounts map[string]uint64) LineEndingType {
	var mostFrequentEnding string
	var highestCount uint64

	for lineEnding, count := range lineEndingCounts {
		if count > highestCount {
			mostFrequentEnding = lineEnding
			highestCount = count
		}
	}
	return LineEndingType(mostFrequentEnding)
}
