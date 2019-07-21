package judge0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// ---------- OU version ----------

// OuSubmit a source code from submission to Judge0 api
func OuSubmit(langID int, source, input, expectedOutput string) *Result {
	var client = &http.Client{}
	
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

	var result *Result
	result = new(Result)
	json.NewDecoder(res.Body).Decode(result)
	return result
}
