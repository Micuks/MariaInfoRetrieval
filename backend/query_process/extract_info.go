package query_process

import (
	. "MariaInfoRetrieval/maria_types"
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func ExtractInfo(doc_id string) (DocumentAbstract, error) {
	log.Info("Extract info for doc ", doc_id)
	var result DocumentAbstract

	doc, ok := idDocMap[doc_id]
	if !ok {
		log.Error("Error getting doc ", doc_id)
	}
	data := map[string]string{"text": doc.Content, "language": doc.Lang.String()}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post(python_server_url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return DocumentAbstract{}, err
	}

	json.NewDecoder(resp.Body).Decode(&result)

	return result, nil
}
