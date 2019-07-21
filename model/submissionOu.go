package model

import (
	"../conversion"
	"../db"
	"../judge0"
)

// OuNewSubmission create a new submission
func OuNewSubmission(userID int, problemID int, languageID int, src string) (int, error) {
	score := 0
	runtime := 0.0
	memory := 0

	testcases, err := OuSpecificTestcaseWithID(problemID)
	if err != nil {
		return -1, err
	}

	ch := make(chan *judge0.Result)

	for i := range testcases {
		go func(testcase *Testcase, i int) {
			ch <- judge0.OuSubmit(languageID, src, testcase.Input, testcase.Output)
		}(testcases[i], i)
	}

	for range testcases {
		result := <- ch
		memory += result.Memory
		runtime += conversion.StringToFloat(result.Time)
		if result.Status.ID == 3 {
			score++
		}
	}

	length := len(testcases)
	runtime = runtime / float64(length)
	memory = memory / length

	

	statement := `INSERT INTO public.submission (usr_id,problem_id,lang_id,src,score,runtime,memory_usage,max_score)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING submission_id`
	row := db.DB.QueryRow(statement, userID, problemID, languageID, src, score, runtime, memory, length)

	lastInsertID := -1
	err = row.Scan(&lastInsertID)
	if err != nil {
		return -1, err
	}

	err = collectScore(userID)
	if err != nil {
		return -1, err
	}

	return lastInsertID, nil
}