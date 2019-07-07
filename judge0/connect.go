package judge0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type request struct {
	source         string
	languageID     int
	stdin          string
	expectedOutput string
}

// Result from Judge0 api
type Result struct {
	Time   string `json:"time"`
	Memory int    `json:"memory"`
	Status struct {
		ID int `json:"id"` // 3 for correct 4 for incorrect
	} `json:"status"`
}

var client = &http.Client{}

// Submit a source code from submission to Judge0 api
func Submit(source, input, expectedOutput string) *Result {
	url := "http://161.246.34.95:4000/submissions?wait=true"

	req := map[string]string{"source_code": source, "language_id": "34", "expected_output": expectedOutput}

	//convert request to io.Reader
	reqByte, _ := json.Marshal(req)
	reqReader := bytes.NewBuffer(reqByte)

	res, err := client.Post(url, "application/json", reqReader)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var result *Result
	result = new(Result)
	json.NewDecoder(res.Body).Decode(result)
	fmt.Println(result)
	return result
}
