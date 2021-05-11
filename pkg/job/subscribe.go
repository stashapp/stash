package job

type ManagerSubscription struct {
	NewJob     <-chan Job
	RemovedJob <-chan Job
	UpdatedJob <-chan Job

	newJob     chan Job
	removedJob chan Job
	updatedJob chan Job
}

func newSubscription() *ManagerSubscription {
	ret := &ManagerSubscription{
		newJob:     make(chan Job, 100),
		removedJob: make(chan Job, 100),
		updatedJob: make(chan Job, 100),
	}

	ret.NewJob = ret.newJob
	ret.RemovedJob = ret.removedJob
	ret.UpdatedJob = ret.updatedJob

	return ret
}

func (s *ManagerSubscription) close() {
	close(s.newJob)
	close(s.removedJob)
	close(s.updatedJob)
}
