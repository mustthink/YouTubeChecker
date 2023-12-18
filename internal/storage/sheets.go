package storage

import (
	"fmt"

	"google.golang.org/api/sheets/v4"

	"github.com/mustthink/YouTubeChecker/internal/types"
)

func (s *VideoStorage) writeToSheet(video types.Video) error {
	values := [][]interface{}{
		video.ArrayInterface(),
	}

	valueRange := &sheets.ValueRange{
		Values: values,
	}

	range_ := fmt.Sprintf("%s!A1", s.config.Name)
	_, err := s.sheets.Spreadsheets.Values.Append(s.config.SpreadsheetID, range_, valueRange).ValueInputOption("RAW").Do()
	return err
}
