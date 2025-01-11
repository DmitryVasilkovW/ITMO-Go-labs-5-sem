package externalsort

import "io"

type WriterImpl struct {
	ioWriter   io.Writer
	lineEnding LineEndingType
}

func (wi *WriterImpl) Write(line string) error {
	_, err := wi.ioWriter.Write([]byte(line))

	if err == nil {
		_, err = wi.ioWriter.Write(wi.lineEnding)
	}

	return err
}

func (wi *WriterImpl) changeLineEnding(lineEnding LineEndingType) *WriterImpl {
	if len(lineEnding) > 0 {
		wi.lineEnding = lineEnding
	}

	return wi
}
