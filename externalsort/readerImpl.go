package externalsort

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type ReaderImpl struct {
	ioReader    io.Reader
	bufioReader *bufio.Reader
	lineEnding  LineEndingType
}

func (ri *ReaderImpl) ReadLine() (string, error) {
	var stringBuilder strings.Builder
	var previousByte byte

	for {
		currentByte, err := ri.readNextByte()
		if err != nil {
			return stringBuilder.String(), err
		}

		if ri.isWindowsLineEnding(previousByte, currentByte) {
			return ri.handleWindowsLineEnding(&stringBuilder), nil
		}

		if ri.isUnixLineEnding(currentByte) {
			return stringBuilder.String(), nil
		}

		stringBuilder.WriteByte(currentByte)
		previousByte = currentByte
	}
}

func (ri *ReaderImpl) readNextByte() (byte, error) {
	return ri.bufioReader.ReadByte()
}

func (ri *ReaderImpl) isWindowsLineEnding(previousByte, currentByte byte) bool {
	return bytes.Equal(ri.lineEnding, LineEnding.Windows) && previousByte == '\r' && currentByte == '\n'
}

func (ri *ReaderImpl) isUnixLineEnding(currentByte byte) bool {
	return currentByte == '\n'
}

func (ri *ReaderImpl) handleWindowsLineEnding(sb *strings.Builder) string {
	result := sb.String()
	return result[:len(result)-1]
}

func (ri *ReaderImpl) changeLineEnding(lineEnding LineEndingType) *ReaderImpl {
	if len(lineEnding) > 0 {
		ri.lineEnding = lineEnding
	}

	return ri
}
