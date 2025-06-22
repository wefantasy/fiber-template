package scheduler

type Task interface {
	Run()
}

func RunTask(tasks []Task) {
	for _, task := range tasks {
		go func() {
			task.Run()
		}()
	}
}

func RunTaskSimple(tasks []Task) {
	for _, task := range tasks {
		task.Run()
	}
}

type CrawlerBase interface {
	GenerateTasks()
	Crawler()
	Filter()
	Persistence()
}

func RunCrawler(c CrawlerBase) {
	c.GenerateTasks()
	c.Crawler()
	c.Filter()
	c.Persistence()
}
