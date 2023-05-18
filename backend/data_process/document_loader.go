package document_process

import (
	. "MariaInfoRetrieval/maria_types"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func LoadDocuments(dir string) ([]Document, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var documents []Document

	for _, file := range files {
		filename := filepath.Join(dir, file.Name())

		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		var doc Document
		dec := json.NewDecoder(f)
		if err := dec.Decode(&doc); err != nil {
			return nil, err
		}

		documents = append(documents, doc)
	}

	return documents, nil
}
