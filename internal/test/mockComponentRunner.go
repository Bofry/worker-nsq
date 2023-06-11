package test

type MockComponentRunner struct {
	prefix string
}

func (c *MockComponentRunner) Start() {
	defaultLogger.Println(c.prefix + ".Start()")
}

func (c *MockComponentRunner) Stop() {
	defaultLogger.Println(c.prefix + ".Stop()")
}
