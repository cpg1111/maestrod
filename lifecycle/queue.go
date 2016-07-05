package lifecycle

import (
	"time"

	"github.com/cpg1111/maestrod/datastore"
)

// QueueEntry is an entry in the waiting and running queue
type QueueEntry struct {
	Project    string
	Branch     string
	PrevCommit string
	CurrCommit string
	CreatedAt  time.Time
	FinishedAt time.Time
	Status     string
}

// Running is the running queue
type Running struct {
	Builds []*QueueEntry
}

// Queue is the qaiting queue
type Queue struct {
	Queue []*QueueEntry
	store *datastore.Datastore
}

// NewQueue returns a pointer to an instanc of a queue
func NewQueue(store *datastore.Datastore) *Queue {
	return &Queue{
		Queue: []*QueueEntry{},
		store: store,
	}
}

func (q *Queue) set(queue []*QueueEntry) {
	q.Queue = queue
}

// Add adds a project to the queue
func (q *Queue) Add(proj, branch, prevCommit, currCommit string) {
	newEntry := &QueueEntry{
		Project:    proj,
		Branch:     branch,
		PrevCommit: prevCommit,
		CurrCommit: currCommit,
		CreatedAt:  time.Now(),
		Status:     "queued",
	}
	q.Queue = append(q.Queue, newEntry)
}

// Pop pops the index 0 of the queue if it can run
func (q *Queue) Pop(r *Running, maxBuilds int) *QueueEntry {
	if len(q.Queue) == 0 {
		return nil
	}
	next := q.Queue[0]
	if len(r.Builds) == 0 {
		r.Builds = []*QueueEntry{next}
		q.Queue = q.Queue[1:]
		return next
	} else if len(r.Builds) < maxBuilds {
		for i := range r.Builds {
			if !(r.Builds[i].Project == next.Project && r.Builds[i].Branch == next.Branch) {
				r.Builds = append(r.Builds, next)
				q.Queue = q.Queue[1:]
				return next
			}
		}
	}
	return nil
}

func (q *Queue) SnapShot() error {
	errChan := make(chan error)
	store := *q.store
	store.Save("queue", q.Queue, func(err error) {
		errChan <- err
	})
	saveErr := <-errChan
	return saveErr
}
