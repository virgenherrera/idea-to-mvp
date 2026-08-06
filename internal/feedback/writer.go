package feedback

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

func Write(feedbackDir string, entry Entry) error {
	if err := os.MkdirAll(feedbackDir, 0755); err != nil {
		return err
	}

	date := time.Now().UTC().Format("2006-01-02")
	filename := filepath.Join(feedbackDir, date+".jsonl")

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	data = append(data, '\n')
	_, err = f.Write(data)
	return err
}
