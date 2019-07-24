package model

import (
	"../db"
)

//Statistic a stat of grader user
type Statistic struct {
	Name          string               `json:"name"`
	Score         float64              `json:"score"`
	ProgressArray []*Progress          `json:"progress_array"`
	Overall       []*OverallSubmission `json:"overall"`
	UserHistory   []*History           `json:"history"`
	ActivePulse   []*Pulse             `json:"active_pulse"`
	Histogram     []*Histogram         `json:"histogram"`
}

// OverallSubmission get the overall submission data of user
type OverallSubmission struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

//History recent submission
type History struct {
	ProblemID	int `json:"problem_id"`
	Title		string `json:"title"`
	Results		string `json:"results"`
	LastDo		string `json:"last_do"`
}

//Pulse Active pulse of user submission
type Pulse struct {
	Date             string `json:"date"`
	NumOfSubmissions int    `json:"submission_count"`
}

//Progress progress of user
type Progress struct {
	CategoryID int `json:"category_id"`
	Completed  int `json:"completed"`
	All        int `json:"all"`
}

//Histogram data for create a histogram
type Histogram struct {
	Start	float32	`json:"start"`
	Stop	float32 `json:"stop"`
	Amount	int 	`json:"amount"`
}

//SpecificUserStatWithID get user stat by id
func SpecificUserStatWithID(id int) (*Statistic, error) {
	statement := `SELECT username,score FROM grader_user WHERE id=$1`
	row := db.DB.QueryRow(statement, id)

	statistic := new(Statistic)
	err := row.Scan(&statistic.Name, &statistic.Score)
	if err != nil {
		return nil, err
	}

	overallSubmissions := make([]*OverallSubmission, 0)

	overall := new(OverallSubmission)
	statement = `SELECT COUNT(distinct problem_id) FROM submission WHERE score=max_score AND usr_id=$1 HAVING COUNT(submission_id)>0;`
	row = db.DB.QueryRow(statement, id)
	row.Scan(&overall.Amount)
	if &overall.Amount == nil {
		overall.Amount = 0
	}
	overall.Name = "completed"
	overallSubmissions = append(overallSubmissions, overall)

	completed := overall.Amount
	var temp int
	overall = new(OverallSubmission)
	statement = `SELECT COUNT(distinct problem_id) FROM submission WHERE usr_id=$1 HAVING COUNT(submission_id) > 0;`
	row = db.DB.QueryRow(statement, id)
	row.Scan(&temp)
	if &temp == nil {
		temp = 0
	}
	overall.Amount = temp - completed
	overall.Name = "working"
	overallSubmissions = append(overallSubmissions, overall)

	overall = new(OverallSubmission)
	statement = `SELECT COUNT(distinct problem_id) FROM submission WHERE usr_id=$1 HAVING COUNT(submission_id) > 0;`
	row = db.DB.QueryRow(statement, id)
	row.Scan(&temp)
	if &temp == nil {
		temp = 0
	}
	overall.Amount = countAllProblems() - temp
	overall.Name = "not started"
	overallSubmissions = append(overallSubmissions, overall)

	statistic.Overall = overallSubmissions

	statement = 
`select
	allProblem.categoryid,
	(case when userStats.count is NULL THEN 0 ELSE userStats.count END) as completed,
	allProblem.count as all
from (
	select a.categoryid, count(a.id)
	from (
		select public.problem.categoryid, public.problem.id
		from public.submission
		inner join public.problem
		on public.submission.problem_id = public.problem.id 
		where public.submission.score = public.submission.max_score and public.submission.usr_id = $1
		group by public.problem.id
	) as a
	group by a.categoryid
) as userStats
full outer join (
	select public.problem.categoryid, count(public.problem.categoryid)
	from public.problem
	group by public.problem.categoryid
) as allProblem
on userStats.categoryid = allProblem.categoryid
order by allProblem.categoryid`
	rows, err := db.DB.Query(statement, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	progressArr := make([]*Progress, 0)
	for rows.Next() {
		progress := new(Progress)

		err := rows.Scan(&progress.CategoryID, &progress.Completed, &progress.All)
		if err != nil {
			return nil, err
		}

		progressArr = append(progressArr, progress)
	}

	statistic.ProgressArray = progressArr

	statement = 
`select submission.problem_id, problem.title, submission.results, submission.last_do
from (
	select distinct on (submission.problem_id) public.submission.problem_id, public.submission.submittedat as last_do, public.submission.results, submission.submission_id
		from public.submission
		where public.submission.usr_id = $1
		and public.submission.problem_id not in (
			select public.submission.problem_id
			from public.submission
			where public.submission.usr_id = $1
			and public.submission.score = public.submission.max_score
			group by public.submission.problem_id
		)
	order by submission.problem_id, submission.submission_id desc
) as submission
inner join (
	select public.problem.id, public.problem.title
	from public.problem
) as problem
on submission.problem_id = problem.id
order by submission.submission_id desc
limit 4`
	rows, err = db.DB.Query(statement, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	histories := make([]*History, 0)
	for rows.Next() {
		history := new(History)

		err := rows.Scan(&history.ProblemID, &history.Title, &history.Results, &history.LastDo)
		if err != nil {
			return nil, err
		}

		histories = append(histories, history)
	}

	statistic.UserHistory = histories

	statement = `SELECT date_trunc('day',submittedat),count(1) FROM submission WHERE usr_id=$1 GROUP BY 1 ORDER BY 1`
	rows, err = db.DB.Query(statement, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pulses := make([]*Pulse, 0)
	for rows.Next() {
		pulse := new(Pulse)

		err := rows.Scan(&pulse.Date, &pulse.NumOfSubmissions)
		if err != nil {
			return nil, err
		}

		pulses = append(pulses, pulse)
	}

	statistic.ActivePulse = pulses

	statement =
`select public.grader_user.score
from public.grader_user
order by public.grader_user.score`
	rows, err = db.DB.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scores := make([]float64, 0)
	for rows.Next() {
		score := 0.0

		err := rows.Scan(&score)
		if err != nil {
			return nil, err
		}

		scores = append(scores, score)
	}

	min := scores[0]
	max := scores[len(scores) - 1]

	histograms := [5]int{0, 0, 0, 0, 0}
	for _, v := range scores {
		histograms[int(float32(v - min) / float32(max - min) * 4)]++
	}

	histograms[3] += histograms[4]
	histogramsSplit := histograms[0:4]

	histogramsOut := make([]*Histogram, 0)

	for k, v := range histogramsSplit {
		histogramOut := new(Histogram)

		histogramOut.Amount = v
		histogramOut.Start = float32(max-min) / 4.0 * float32(k)
		histogramOut.Stop = float32(max-min) / 4.0 * float32(k + 1)

		histogramsOut = append(histogramsOut, histogramOut)
	}

	statistic.Histogram = histogramsOut

	return statistic, nil
}
