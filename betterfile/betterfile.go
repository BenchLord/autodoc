package betterfile

import "os"

type BetterFile struct {
	*os.File
}

func NewFile(name string) (*BetterFile, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return &BetterFile{f}, nil
}

func (f *BetterFile) GetLineAsBytes() ([]byte, error) {
	line := make([]byte, 0)
	for {
		buffer := make([]byte, 1)
		_, err := f.Read(buffer)
		if err != nil {
			return nil, err
		}
		line = append(line, buffer[0])
		if buffer[0] == 10 {
			break
		}
	}
	return line, nil
}
