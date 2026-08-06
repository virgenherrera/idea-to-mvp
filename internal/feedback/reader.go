package feedback

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func ReadAll(feedbackDir string) ([]Entry, error) {
	files, err := filepath.Glob(filepath.Join(feedbackDir, "*.jsonl"))
	if err != nil {
		return nil, err
	}

	var entries []Entry

	for _, file := range files {
		fileEntries, err := readFile(file)
		if err != nil {
			return nil, err
		}
		entries = append(entries, fileEntries...)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp > entries[j].Timestamp
	})

	return entries, nil
}

func readFile(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var entry Entry
		if err := json.Unmarshal(line, &entry); err != nil {
			fmt.Fprintf(os.Stderr, "warning: skipping malformed line %d in %s: %v\n", lineNum, path, err)
			continue
		}
		entries = append(entries, entry)
	}

	return entries, scanner.Err()
}
