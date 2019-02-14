package manager

type JobStatus int

const (
	Idle     JobStatus = 0
	Import   JobStatus = 1
	Export   JobStatus = 2
	Scan     JobStatus = 3
	Generate JobStatus = 4
	Clean    JobStatus = 5
	Scrape   JobStatus = 6
)
