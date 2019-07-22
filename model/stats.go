package model

import (
	"ce-boostup-backend/db"
	"fmt"
)

//Statistic a stat of grader user
type Statistic struct {
	Name          string               `json:"name"`
	Score         int                  `json:"score"`
	ProgressArray []*Progress          `json:"progress_array"`
	Overall       []*OverallSubmission `json:"overall"`
	UserHistory   []*History           `json:"history"`
	ActivePulse   []*Pulse             `json:"active_pulse"`
	Histogram     Histogram            `json:"histogram"`
}

// OverallSubmission get the overall submission data of user
type OverallSubmission struct {
	Name   string  `json:"name"`
	Amount float64 `json:"value"`
}

//History recent submission
type History struct {
	ProblemID	int `json:"problem_id"`
	Title		string `json:"title"`
	Percen		float32 `json:"percen"`
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
	AllPlot   []*Plot `json:"plot"`
	UserScore int     `json:"user_score"`
}

//Plot plot a histogram
type Plot struct {
	ProgressRange string `json:"progress_range"`
	Amount        int    `json:"amount"`
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
	completed := overall.Amount
	overall.Amount = (overall.Amount / float64(countAllProblems())) * 100.0
	overall.Name = "success"
	overallSubmissions = append(overallSubmissions, overall)

	var temp int
	overall = new(OverallSubmission)
	statement = `SELECT COUNT(distinct problem_id) FROM submission WHERE usr_id=$1 HAVING COUNT(submission_id) > 0;`
	row = db.DB.QueryRow(statement, id)
	row.Scan(&temp)
	if &temp == nil {
		temp = 0
	}
	fmt.Println(temp, completed)
	overall.Amount = (float64(temp) - completed) / float64(countAllProblems()) * 100.0
	overall.Name = "in progress"
	overallSubmissions = append(overallSubmissions, overall)

	overall = new(OverallSubmission)
	statement = `SELECT COUNT(distinct problem_id) FROM submission WHERE usr_id=$1 HAVING COUNT(submission_id) > 0;`
	row = db.DB.QueryRow(statement, id)
	row.Scan(&temp)
	if &temp == nil {
		temp = 0
	}
	overall.Amount = (float64(countAllProblems()) - float64(temp)) / float64(countAllProblems()) * 100.0
	overall.Name = "not started"
	overallSubmissions = append(overallSubmissions, overall)

	statistic.Overall = overallSubmissions

	statement = 
`select
	allProblem.categoryid,
	(case when userStats.count is NULL THEN 0 ELSE userStats.count END) as completed,
	allProblem.count as all
from (select public.problem.categoryid, count(public.problem.categoryid)
	from public.submission
	inner join public.problem
	on public.submission.problem_id = public.problem.id 
	where public.submission.score = public.submission.max_score and public.submission.usr_id = $1
	group by public.problem.categoryid) as userStats
full outer join (select public.problem.categoryid, count(public.problem.categoryid)
	from public.problem
	group by public.problem.categoryid) as allProblem
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
`select submission.problem_id, problem.title, submission.percen, submission.last_do
from (
	select public.submission.problem_id,
		max(public.submission.submittedat) as last_do,
		max(public.submission.score * 100.0 / public.submission.max_score) as percen
	from public.submission
	where public.submission.usr_id = $1
		and public.submission.problem_id not in (
			select public.submission.problem_id 
			from public.submission
			where public.submission.score = public.submission.max_score
			group by public.submission.problem_id
		)
	group by public.submission.problem_id
) as submission
inner join (
	select public.problem.id, public.problem.title
	from public.problem
) as problem
on submission.problem_id = problem.id
order by submission.last_do desc
limit 4`
	rows, err = db.DB.Query(statement, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	histories := make([]*History, 0)
	for rows.Next() {
		history := new(History)

		err := rows.Scan(&history.ProblemID, &history.Title, &history.Percen, &history.LastDo)
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

	statement = `SELECT range,freq FROM histogram('grader_user','score')`
	rows, err = db.DB.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plots := make([]*Plot, 0)
	for rows.Next() {
		plot := new(Plot)

		err := rows.Scan(&plot.ProgressRange, &plot.Amount)
		if err != nil {
			return nil, err
		}

		plots = append(plots, plot)
	}

	statistic.Histogram.UserScore = statistic.Score
	statistic.Histogram.AllPlot = plots

	return statistic, nil
}
