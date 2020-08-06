package manager

// PostMigrate is executed after migrations have been executed.
func (s *singleton) PostMigrate() {
	setInitialMD5Config()
}
