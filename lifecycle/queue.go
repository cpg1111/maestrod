package lifecycle

import (
	"time"

	"github.com/cpg1111/maestrod/datastore"
)

type QueueEntry struct {
	Project    string
	Branch     string
	Commit     string
	ConfDir    string
	CreatedAt  time.Time
	FinishedAt time.Time
	Status     string
}

type Running struct {
	Builds []*QueueEntry
}

type Queue struct {
	Queue []*QueueEntry
	store *datastore.Datastore
}

func NewQueue(store *datastore.Datastore) *Queue {
	return &Queue{
		Queue: []*QueueEntry{},
		store: store,
	}
}

func (q *Queue) set(queue []*QueueEntry) {
	q.Queue = queue
}

func (q *Queue) Add(proj, branch string) {
	newEntry := &QueueEntry{
		Project:   proj,
		Branch:    branch,
		CreatedAt: time.Now(),
		Status:    "queued",
	}
	q.Queue = append(q.Queue, newEntry)
}

func (q *Queue) Pop(r *Running, maxBuilds int) *QueueEntry {
	if len(q.Queue) == 0 {
		return nil
	}
	next := q.Queue[0]
	if len(r.Builds) == 0 {
		r.Builds = []*QueueEntry{next}
		q.Queue = q.Queue[1:]
	} else if len(r.Builds) < maxBuilds {
		for i := range r.Builds {
			if !(r.Builds[i].Project == next.Project && r.Builds[i].Branch == next.Branch) {
				r.Builds = append(r.Builds, next)
				q.Queue = q.Queue[1:]
			}
		}
	}
	return nil
}
