package judge0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

// ---------- OU version ----------

// OuResult from Judge0 api
type OuResult struct {
	Time			string 	`json:"time"`
	Memory			int    	`json:"memory"`
	CompileOutput	string	`json:"compile_output"`
	Status struct {
		ID int `json:"id"` // 3 for correct 4 for incorrect
	} `json:"status"`
}

// OuSubmit a source code from submission to Judge0 api
func OuSubmit(langID int, source, input, expectedOutput string) *OuResult {
	url := fmt.Sprintf("http://%s:%s/submissions?wait=true", os.Getenv("JUDGE_0_IP"), os.Getenv("JUDGE_0_PORT"))

	req := map[string]string{"source_code": source, "stdin": input, "language_id": strconv.Itoa(langID), "expected_output": expectedOutput}

	//convert request to io.Reader
	reqByte, _ := json.Marshal(req)
	reqReader := bytes.NewBuffer(reqByte)

	res, err := client.Post(url, "application/json", reqReader)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	result := new(OuResult)
	json.NewDecoder(res.Body).Decode(result)
	return result
}
