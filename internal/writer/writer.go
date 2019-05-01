package writer

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/sceptero/house-extractor/internal/types"
)

// File -
type File struct {
	file *os.File
}

// New File
func New(filePath string) (*File, error) {
	f, err := os.Create(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "while opening file for writing")
	}
	return &File{file: f}, nil
}

// CloseFile -
func (f *File) CloseFile() error {
	return f.file.Close()
}

func (f *File) Write(houses map[int][]types.HouseTile) error {
	f.file.WriteString(fmt.Sprintf("houses = {\n"))
	for id, tiles := range houses {
		f.file.WriteString(fmt.Sprintf("  [%d] = {\n", id))
		for _, tile := range tiles {
			f.file.WriteString(fmt.Sprintf("    {x = %d, y = %d, z = %d},\n", tile.PosX, tile.PosY, tile.PosZ))
		}
		f.file.WriteString(fmt.Sprintf("  },\n"))
	}
	f.file.WriteString(fmt.Sprintf("}\n"))
	return nil
}
