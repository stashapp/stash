package job

func ExecuteTask(p *Progress, description string, t func()) {
	task := &Task{
		description: description,
	}

	p.AddTask(task)
	defer p.RemoveTask(task)
	t()
}
