package cleaner

import "os"

var works []func()

func Defer(work func()) {
	works = append(works, work)
}

func Work() {
	for _, work := range works {
		work()
	}
}

func Exit(status int) {
	Work()
	os.Exit(status)
}
