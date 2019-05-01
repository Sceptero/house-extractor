package extractor

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"

	"github.com/sceptero/house-extractor/internal/reader"
	"github.com/sceptero/house-extractor/internal/types"
	"github.com/sceptero/house-extractor/internal/writer"
)

// HouseExtractor -
type HouseExtractor struct {
	inputFile  *reader.File
	outputFile *writer.File
	houses     map[int][]types.HouseTile
	tilesCount int
}

// New HouseExtractor
func New(inputFilePath, outputFilePath string) (*HouseExtractor, error) {
	iFile, err := reader.New(inputFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "while creating file reader")
	}

	oFile, err := writer.New(outputFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "while creating file writer")
	}

	houses := make(map[int][]types.HouseTile)

	return &HouseExtractor{inputFile: iFile, outputFile: oFile, tilesCount: 0, houses: houses}, nil
}

// Do house data extraction
func (e *HouseExtractor) Do() error {
	err := e.readMapData()
	if err != nil {
		return errors.Wrap(err, "while reading map data")
	}
	e.inputFile.CloseFile()

	err = e.outputFile.Write(e.houses)
	if err != nil {
		return errors.Wrap(err, "while writing house data to file")
	}
	e.outputFile.CloseFile()

	return nil
}

func (e *HouseExtractor) readMapData() error {
	err := e.validateIdentifier()
	if err != nil {
		return errors.Wrap(err, "while validating file identifier")
	}

	// look for tile area
	for {
		_, err = e.inputFile.SeekBytesWithTerminator([]byte{reader.NodeStart, reader.TileArea}, nil)
		if err != nil {
			if errors.Cause(err).Error() == "EOF" {
				break
			}
			return errors.Wrap(err, "while looking for TileArea node in file")
		}

		tileArea, err := e.readTileArea()
		if err != nil {
			return errors.Wrap(err, "while reading tile area")
		}

		// look for house tile
		for {
			terminated, err := e.inputFile.SeekBytesWithTerminator([]byte{reader.NodeStart, reader.HouseTile}, []byte{reader.NodeStart, reader.TileArea})
			if err != nil {
				if errors.Cause(err).Error() == "EOF" {
					break
				}
				return errors.Wrap(err, "while looking for byte sequence in file")
			}
			if terminated {
				break
			}

			houseTile, err := e.readHouseTile(tileArea)
			if err != nil {
				return errors.Wrap(err, "while reading house tile")
			}
			e.houses[houseTile.ID] = append(e.houses[houseTile.ID], houseTile)
			e.tilesCount++
		}
	}

	fmt.Printf("Houses found: %v, HouseTiles found: %v\n", len(e.houses), e.tilesCount)
	return nil
}

func (e *HouseExtractor) readHouseTile(tileArea types.TileArea) (types.HouseTile, error) {
	offX, err := e.inputFile.ReadU8()
	if err != nil {
		return types.HouseTile{}, errors.Wrap(err, "while reading house tile offset x pos")
	}
	offY, err := e.inputFile.ReadU8()
	if err != nil {
		return types.HouseTile{}, errors.Wrap(err, "while reading house tile offset y pos")
	}
	houseID, err := e.inputFile.ReadU32()
	if err != nil {
		return types.HouseTile{}, errors.Wrap(err, "while reading house tile house id")
	}

	return types.HouseTile{
		PosX: (tileArea.BaseX + int(offX)),
		PosY: (tileArea.BaseY + int(offY)),
		PosZ: tileArea.BaseZ,
		ID:   int(houseID),
	}, nil
}

func (e *HouseExtractor) readTileArea() (types.TileArea, error) {
	baseX, err := e.inputFile.ReadU16()
	if err != nil {
		return types.TileArea{}, errors.Wrap(err, "while reading tile area base x pos")
	}
	baseY, err := e.inputFile.ReadU16()
	if err != nil {
		return types.TileArea{}, errors.Wrap(err, "while reading tile area base y pos")
	}
	baseZ, err := e.inputFile.ReadU8()
	if err != nil {
		return types.TileArea{}, errors.Wrap(err, "while reading tile area base z pos")
	}

	return types.TileArea{
		BaseX: int(baseX),
		BaseY: int(baseY),
		BaseZ: int(baseZ),
	}, nil
}

func (e *HouseExtractor) validateIdentifier() error {
	identifier, err := e.inputFile.ReadBytes(4)
	if err != nil {
		return errors.Wrap(err, "while reading file identifier")
	}

	if !bytes.Equal([]byte{0, 0, 0, 0}, identifier) && string(identifier) != "OTBM" {
		return errors.New(fmt.Sprintf("invalid file identifier: %s", identifier))
	}

	return nil
}
