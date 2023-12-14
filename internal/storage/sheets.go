package storage

import (
	"google.golang.org/api/sheets/v4"
)

func (s *VideoStorage) writeToSheet(video Video) error {
	values := [][]interface{}{
		{video.videoID, video.videoTitle, video.channelID, video.channelTitle, video.description, video.publishDate, video.videoURL, video.thumbnailURL},
	}

	valueRange := &sheets.ValueRange{
		Values: values,
	}

	_, err := s.sheets.Spreadsheets.Values.Append(s.config.SpreadsheetID, "A1", valueRange).ValueInputOption("RAW").Do()
	return err
}
