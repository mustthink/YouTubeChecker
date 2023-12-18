package storage

import (
	"fmt"

	"google.golang.org/api/sheets/v4"
)

func (s *VideoStorage) writeToSheet(video Video) error {
	values := [][]interface{}{
		{video.isTracked, video.videoID, video.videoTitle, video.channelID, video.channelTitle, video.description, video.publishDate, video.videoURL, video.thumbnailURL},
	}

	valueRange := &sheets.ValueRange{
		Values: values,
	}

	range_ := fmt.Sprintf("%s!A1", s.config.Name)
	_, err := s.sheets.Spreadsheets.Values.Append(s.config.SpreadsheetID, range_, valueRange).ValueInputOption("RAW").Do()
	return err
}
