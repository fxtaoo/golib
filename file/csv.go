package file

import (
	"encoding/csv"
	"io"
	"os"
)

// CSV 文件读数据
func CSVRead(filePath string) ([][]string, error) {
	filePath = AbsPath(filePath)

	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.LazyQuotes = true

	var csvdata [][]string
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		// 追加
		csvdata = append(csvdata, row)
	}
	return csvdata, nil
}
