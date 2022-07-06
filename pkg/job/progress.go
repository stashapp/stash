package job

import "sync"

// ProgressIndefinite is the special percent value to indicate that the
// percent progress is not known.
const ProgressIndefinite float64 = -1

// Progress is used by JobExec to communicate updates to the job's progress to
// the JobManager.
type Progress struct {
	defined      bool
	processed    int
	total        int
	percent      float64
	currentTasks []*task

	mutex   sync.Mutex
	updater *updater
}

type task struct {
	description string
}

func (p *Progress) updated() {
	var details []string
	for _, t := range p.currentTasks {
		details = append(details, t.description)
	}

	p.updater.updateProgress(p.percent, details)
}

// Indefinite sets the progress to an indefinite amount.
func (p *Progress) Indefinite() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.defined = false
	p.total = 0
	p.calculatePercent()
}

// Definite notifies that the total is known.
func (p *Progress) Definite() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.defined = true
	p.calculatePercent()
}

// SetTotal sets the total number of work units and sets definite to true.
// This is used to calculate the progress percentage.
func (p *Progress) SetTotal(total int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.total = total
	p.defined = true
	p.calculatePercent()
}

// AddTotal adds to the total number of work units. This is used to calculate the
// progress percentage.
func (p *Progress) AddTotal(total int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.total += total
	p.calculatePercent()
}

// SetProcessed sets the number of work units completed. This is used to
// calculate the progress percentage.
func (p *Progress) SetProcessed(processed int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.processed = processed
	p.calculatePercent()
}

func (p *Progress) calculatePercent() {
	switch {
	case !p.defined || p.total <= 0:
		p.percent = ProgressIndefinite
	case p.processed < 0:
		p.percent = 0
	default:
		p.percent = float64(p.processed) / float64(p.total)
		if p.percent > 1 {
			p.percent = 1
		}
	}

	p.updated()
}

// SetPercent sets the progress percent directly. This value will be
// overwritten if Indefinite, SetTotal, Increment or SetProcessed is called.
// Constrains the percent value between 0 and 1, inclusive.
func (p *Progress) SetPercent(percent float64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if percent < 0 {
		percent = 0
	} else if percent > 1 {
		percent = 1
	}

	p.percent = percent
	p.updated()
}

// Increment increments the number of processed work units. This is used to calculate the percentage.
// If total is set already, then the number of processed work units will not exceed the total.
func (p *Progress) Increment() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.defined || p.total <= 0 || p.processed < p.total {
		p.processed++
		p.calculatePercent()
	}
}

// AddProcessed increments the number of processed work units by the provided
// amount. This is used to calculate the percentage.
func (p *Progress) AddProcessed(v int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	newVal := v
	if p.defined && p.total > 0 && newVal > p.total {
		newVal = p.total
	}

	p.processed = newVal
	p.calculatePercent()
}

func (p *Progress) addTask(t *task) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.currentTasks = append([]*task{t}, p.currentTasks...)
	p.updated()
}

func (p *Progress) removeTask(t *task) {
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

// ExecuteTask executes a task as part of a job. The description is used to
// populate the Details slice in the parent Job.
func (p *Progress) ExecuteTask(description string, fn func()) {
	t := &task{
		description: description,
	}

	p.addTask(t)
	defer p.removeTask(t)
	fn()
}
