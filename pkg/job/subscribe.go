package job

// ManagerSubscription is a collection of channels that will receive updates
// from the job manager.
type ManagerSubscription struct {
	// new jobs are sent to this channel
	NewJob <-chan Job
	// removed jobs are sent to this channel
	RemovedJob <-chan Job
	// updated jobs are sent to this channel
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
