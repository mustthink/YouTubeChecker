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

	range_ := fmt.Sprintf("%s!A2", s.config.Name)
	_, err := s.sheets.Spreadsheets.Values.Append(s.config.SpreadsheetID, range_, valueRange).ValueInputOption("RAW").Do()
	return err
}

func (s *VideoStorage) SortSheet() error {
	sortRange := &sheets.GridRange{
		SheetId:          s.config.SheetID,
		StartRowIndex:    1,
		StartColumnIndex: 0,
		EndRowIndex:      0,
		EndColumnIndex:   0,
	}

	sortRequest := sheets.SortRangeRequest{
		Range:     sortRange,
		SortSpecs: []*sheets.SortSpec{{DimensionIndex: 8, SortOrder: "DESCENDING"}},
	}

	reqs := []*sheets.Request{
		{
			SortRange: &sortRequest,
		},
	}

	batchUpdateReq := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: reqs,
	}

	_, err := s.sheets.Spreadsheets.BatchUpdate(s.config.SpreadsheetID, batchUpdateReq).Do()
	return err
}
