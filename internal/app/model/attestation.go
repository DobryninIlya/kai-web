package model

type Discipline struct {
	Number           int
	Name             string
	Assessments      []Assessment
	PreliminaryGrade string
	AdditionalPoints int
	Debts            int
	FinalGrade       int
	TraditionalGrade string
}

type Assessment struct {
	YourScore int
	MaxScore  int
}
