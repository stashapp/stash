package job

import "sync"

// Progress is used by JobExec to communicate updates to the job's progress to
// the JobManager.
type Progress struct {
	processed    int
	total        int
	percent      float64
	currentTasks []*Task

	mutex   sync.Mutex
	updater *updater
}

type Task struct {
	description string
}

func (p *Progress) updated() {
	var details []string
	for _, t := range p.currentTasks {
		details = append(details, t.description)
	}

	p.updater.UpdateProgress(p.percent, details)
}

func (p *Progress) Indefinite() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.total = -1
	p.calculatePercent()
}

func (p *Progress) SetTotal(total int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.total = total
	p.calculatePercent()
}

func (p *Progress) SetProcessed(processed int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.processed = processed
	p.calculatePercent()
}

func (p *Progress) calculatePercent() {
	if p.total <= 0 {
		p.percent = -1
	} else {
		p.percent = float64(p.processed) / float64(p.total)
		if p.percent > 1 {
			p.percent = 1
		}
	}

	p.updated()
}

func (p *Progress) SetPercent(percent float64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.percent = percent
	p.updated()
}

func (p *Progress) Increment() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.total > 0 {
		p.processed += 1
		p.calculatePercent()
	}
}

func (p *Progress) AddTask(t *Task) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.currentTasks = append(p.currentTasks, t)
	p.updated()
}

func (p *Progress) RemoveTask(t *Task) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for i, tt := range p.currentTasks {
		if tt == t {
			p.currentTasks = append(p.currentTasks[:i], p.currentTasks[i+1:]...)
			p.updated()
			return
		}
	}
}
