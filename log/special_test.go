package log

import (
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/arutselvan15/go-utils/logconstants"
)

func doSomething(t *testing.T, log *Log) {
	process := "DispositionProblem"
	log.SetDisposition(logconstants.Success)

	logAndAssertJSON(t, log, "Processing", func(fields logrus.Fields) {
		assert.Equal(t, process, fields["process"])
		assert.Equal(t, logconstants.Start, fields["action"])
		assert.Equal(t, logconstants.Success, fields["disposition"])
	})
}

func somePeriodicFunction(t *testing.T, log *Log, counter int) {
	process := "DispositionProblem"
	fmt.Printf("Tick %d\n", counter)
	msg := fmt.Sprintf("Processing %d", counter)

	// At the start of this function, disposition should be empty
	fmt.Println("Before push")
	fmt.Printf("   Disposition = %s\n", log.disposition)
	assert.Empty(t, log.disposition, msg)

	log.PushContext()
	log.SetAction(logconstants.Start)
	log.SetInvolvedObj("SomeKey")

	fmt.Println("After push, doing log test")
	fmt.Printf("   Disposition = %s\n", log.disposition)
	logAndAssertJSON(t, log, msg, func(fields logrus.Fields) {
		assert.Equal(t, process, fields["process"])
		assert.Equal(t, logconstants.Start, fields["action"])
		assert.Equal(t, logconstants.Start, fields["action"])
	})

	fmt.Println("Do something")
	fmt.Printf("   Disposition = %s\n", log.disposition)
	// this will set the disposition
	doSomething(t, log)
	assert.NotEmpty(t, log.disposition, msg)

	fmt.Println("Before Pop")
	fmt.Printf("   Disposition = %s\n", log.disposition)
	log.PopContext()
	assert.Empty(t, log.disposition, msg)
	fmt.Println("After Pop")
	fmt.Printf("   Disposition = %s\n", log.disposition)
}

// Special test for a certain bug condition
func TestLog_DispositionProblem(t *testing.T) {
	log := newLogger()

	process := "DispositionProblem"
	log.SetProcess(process)
	//setAllFields(log, "initial")
	//assertAllFields(t, log, "initial")

	threadLogger := log.ThreadLogger().SetSubProcess("periodic-process")

	ticker := time.NewTicker(1 * time.Second)
	stopchan := make(chan interface{})

	f := func() {
		counter := 0
		for {
			select {
			case <-ticker.C:
				somePeriodicFunction(t, threadLogger, counter)
				counter = counter + 1
			case <-stopchan:
				return
			}
		}
	}

	go f()
	time.Sleep(2 * time.Second)
	stopchan <- true

	log.PopContext()
}
