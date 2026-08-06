package feedback

import "encoding/json"

func Export(entries []Entry) ([]byte, error) {
	return json.MarshalIndent(entries, "", "  ")
}
