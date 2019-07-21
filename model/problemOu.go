package model

import (
	"../db"
)

// ---------- OU version ----------

// OuNewTestcase add a new testcase with just input and we'll get output of testcase by using Judge0
func OuNewTestcase(id int, testcase Testcase) error {
	statement := `UPDATE problem SET testcase = array_append(testcase, ($1,$2)::TESTCASE) WHERE id=$3`
	_, err := db.DB.Exec(statement, testcase.Input, testcase.Output, id)
	if err != nil {
		return err
	}
	return nil
}

// OuSpecificTestcaseWithID return turncase from specific problem
func OuSpecificTestcaseWithID(id int) ([]*Testcase, error) {
	statement := `select (m.u).input, (m.u).output from (SELECT UNNEST(testcase) as u FROM public.problem WHERE id=$1) as m`
	rows, err := db.DB.Query(statement, id)
	if err != nil {
		return nil, err
	}

	testcases := make([]*Testcase, 0)
	for rows.Next() {
		testcase := new(Testcase)
		err := rows.Scan(&testcase.Input, &testcase.Output)
		if err != nil {
			return nil, err
		}
		testcases = append(testcases, testcase)
	}

	return testcases, nil
}

// OuUpdateTestcase add a new testcase with just input and we'll get output of testcase by using Judge0
func OuUpdateTestcase(id int, index int, testcase Testcase) error {
	statement := `UPDATE public.problem SET testcase = (testcase[0:$1] || ($3,$4)::TESTCASE || testcase[$2:]) WHERE id=$5`
	_, err := db.DB.Exec(statement, index - 1, index + 1, testcase.Input, testcase.Output, id)
	if err != nil {
		return err
	}
	return nil
}

// OuDeleteTestcase add a new testcase with just input and we'll get output of testcase by using Judge0
func OuDeleteTestcase(id int, index int) error {
	statement := `UPDATE public.problem SET testcase = (testcase[0:$1]  || testcase[$2:]) WHERE id=$3`
	_, err := db.DB.Exec(statement, index - 1, index + 1, id)
	if err != nil {
		return err
	}
	return nil
}