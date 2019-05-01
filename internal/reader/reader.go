package reader

import (
	"encoding/binary"
	"os"

	"github.com/pkg/errors"
)

const (
	NodeStart  = 0xfe
	NodeEnd    = 0xff
	EscapeChar = 0xfd
	HouseTile  = 0x0e
	TileArea   = 0x04
)

//
// https://github.com/hjnilsson/rme/blob/master/source/iomap_otbm.cpp#L748
// https://github.com/hjnilsson/rme/blob/master/source/iomap_otbm.h#L56
// https://github.com/edubart/otclient/blob/master/src/client/mapio.cpp#L109
//

// File -
type File struct {
	file *os.File
}

// New File
func New(filePath string) (*File, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "while opening file")
	}
	return &File{file: f}, nil
}

// ReadBytes -
func (f *File) ReadBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := f.file.Read(bytes)
	if err != nil {
		return nil, errors.Wrap(err, "while reading bytes from file")
	}
	return bytes, nil
}

// ReadU8 -
func (f *File) ReadU8() (uint8, error) {
	buf, err := f.ReadBytes(1)
	if err != nil {
		return 0, errors.Wrap(err, "while reading u8")
	}
	return buf[0], nil
}

// ReadU16 -
func (f *File) ReadU16() (uint16, error) {
	buf, err := f.ReadBytes(2)
	if err != nil {
		return 0, errors.Wrap(err, "while reading u16")
	}
	return binary.LittleEndian.Uint16(buf), nil
}

// ReadU32 -
func (f *File) ReadU32() (uint32, error) {
	buf, err := f.ReadBytes(4)
	if err != nil {
		return 0, errors.Wrap(err, "while reading u32")
	}
	return binary.LittleEndian.Uint32(buf), nil
}

// ReadU64 -
func (f *File) ReadU64() (uint64, error) {
	buf, err := f.ReadBytes(8)
	if err != nil {
		return 0, errors.Wrap(err, "while reading u64")
	}
	return binary.LittleEndian.Uint64(buf), nil
}

// Skip -
func (f *File) Skip(n int64) error {
	_, err := f.file.Seek(n, 1)
	return err
}

// CloseFile -
func (f *File) CloseFile() error {
	return f.file.Close()
}

// SeekBytesWithTerminator -
func (f *File) SeekBytesWithTerminator(seeked, terminator []byte) (bool, error) {
	if seeked == nil || len(seeked) == 0 {
		return false, nil
	}

	for {
		b, err := f.ReadBytes(1)
		if err != nil {
			return false, errors.Wrap(err, "while reading byte from file")
		}

		// check if terminator
		if terminator != nil && len(terminator) != 0 && b[0] == terminator[0] {
			for id, t := range terminator[1:] {
				b, err := f.ReadBytes(1)
				if err != nil {
					return false, errors.Wrap(err, "while reading byte from file")
				}

				// chain broken, move file offset back
				if b[0] != t {
					_, err = f.file.Seek(int64(-1*(id+1)), 1)
					if err != nil {
						return false, errors.Wrap(err, "while moving file offset back")
					}
					break
				}

				// found full chain
				if id == (len(terminator) - 2) {
					_, err = f.file.Seek(int64(-len(terminator)), 1)
					if err != nil {
						return false, errors.Wrap(err, "while moving file offset back")
					}
					return true, nil
				}
			}
		}

		// found beginning of chain
		if b[0] == seeked[0] {
			for id, s := range seeked[1:] {
				b, err := f.ReadBytes(1)
				if err != nil {
					return false, errors.Wrap(err, "while reading byte from file")
				}

				// chain broken
				if b[0] != s {
					break
				}

				// found full chain
				if id == (len(seeked) - 2) {
					return false, nil
				}
			}
		}
	}
}
