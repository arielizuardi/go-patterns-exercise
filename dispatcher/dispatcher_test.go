package dispatcher_test

import (
	"strconv"
	"testing"

	"github.com/arielizuardi/go-patterns-exercise/dispatcher"
)

func TestSend(t *testing.T) {

	d := dispatcher.NewDispatcher(10)
	go d.Run()

	dispatcher.JobQueue = make(chan dispatcher.Job)
	for i := 1; i < 100; i++ {
		dispatcher.JobQueue <- dispatcher.Job{DeviceToken: `dvc` + strconv.Itoa(i), NotificationID: `push_id`}
	}
}
