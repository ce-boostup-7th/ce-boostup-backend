package judge0

import (
	"strings"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type request struct {
	source         string
	languageID     int
	stdin          string
	expectedOutput string
}

// Result from Judge0 api Ou
type Result struct {
	Time			string 	`json:"time"`
	Memory			int    	`json:"memory"`
	CompileOutput	string	`json:"compile_output"`
	Stdout			*string `json:"stdout"`
	Status struct {
		ID int `json:"id"` // 3 for correct 4 for incorrect
	} `json:"status"`
}

var client = &http.Client{}

// Submit a source code from submission to Judge0 api Ou
func Submit(langID int, source, input, expectedOutput string) *Result {
	url := fmt.Sprintf("http://%s:%s/submissions?wait=true", os.Getenv("JUDGE_0_IP"), os.Getenv("JUDGE_0_PORT"))

	req := map[string]string{"source_code": source, "stdin": input, "language_id": strconv.Itoa(langID)}

	if expectedOutput != "" {
		req["expected_output"] = expectedOutput
	}

	//convert request to io.Reader
	reqByte, _ := json.Marshal(req)
	reqReader := bytes.NewBuffer(reqByte)

	res, err := client.Post(url, "application/json", reqReader)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	result := new(Result)
	json.NewDecoder(res.Body).Decode(result)
	if expectedOutput == "" {
		if result.Stdout == nil {
			result.Stdout = new(string)
		}
		if strings.TrimSpace(*result.Stdout) != "" {
			result.Status.ID = 4
		}
	}
	return result
}
