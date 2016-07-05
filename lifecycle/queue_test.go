package lifecycle

import (
	"os"
	"testing"

	"github.com/cpg1111/maestrod/datastore"
)

var store = datastore.NewRedis(os.Getenv("REDIS_SERVICE_HOST"), os.Getenv("REDIS_SERVICE_PORT"), "")

var castedStore = (datastore.Datastore)(*store)

var queue = NewQueue(&castedStore)

func TestAdd(t *testing.T) {
	queue.Add("test", "test", "asdfasdfasdf", "asdfew")
	if !(queue.Queue[0].Project == "test" && queue.Queue[0].Branch == "test" && queue.Queue[0].PrevCommit == "asdfasdfasdf" && queue.Queue[0].CurrCommit == "asdfew") {
		t.Error("Add did not add all the fields")
	}
}

func TestPop(t *testing.T) {
	queue.Add("test", "test", "asdfasdfasdf", "asdfew")
	running := &Running{}
	next := queue.Pop(running, 1)
	if next == nil {
		t.Error("queue should've popped an entry")
	}
	another := queue.Pop(running, 1)
	if another != nil {
		t.Error("queue's second pop should've returned nil")
	}
}
