/*
Copyright 2016 Christian Grabowski All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lifecycle

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cpg1111/maestrod/datastore"
	"github.com/cpg1111/maestrod/manager"
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

type aliveKey struct {
	Project string
	Branch  string
}

// Running is the running queue
type Running struct {
	Builds    []*QueueEntry
	Alive     map[aliveKey]*QueueEntry
	KeepAlive map[aliveKey]bool
}

// CheckIn changes whether a build is alive or not
func (r *Running) CheckIn(project, branch string) {
	key := aliveKey{
		Project: project,
		Branch:  branch,
	}
	if r.Alive[key] != nil {
		r.KeepAlive[key] = true
	}
}

// Watch watches running maestro builds
func (r *Running) Watch(manager *manager.Driver) {
	for b := range r.Builds {
		key := aliveKey{
			Project: r.Builds[b].Project,
			Branch:  r.Builds[b].Branch,
		}
		if r.Alive[key] == nil {
			r.Builds = append(r.Builds[:b], r.Builds[b+1:]...)
			break
		}
		if !r.KeepAlive[key] {
			r.Alive[key] = nil
			mgrDRef := *manager
			mgrDRef.DestroyWorker(key.Project, key.Branch)
		}
	}
}

// Queue is the qaiting queue
type Queue struct {
	Queue []*QueueEntry
	store datastore.Datastore
}

// NewQueue returns a pointer to an instanc of a queue
func NewQueue(store datastore.Datastore) *Queue {
	return &Queue{
		Queue: []*QueueEntry{},
		store: store,
	}
}

func (q *Queue) set(queue []*QueueEntry) {
	q.Queue = queue
}

type lastSuccess struct {
	Commit string
}

func successKey(proj, branch string) string {
	if len(branch) > 0 {
		return fmt.Sprintf("/success/%s/%s", proj, branch)
	}
	return fmt.Sprintf("/success/%s", proj)
}

// GetLastSuccess fetches the commit hash of the last successful build
func (q *Queue) GetLastSuccess(proj, branch string) (string, error) {
	resChan := make(chan string)
	errChan := make(chan error)
	key := successKey(proj, branch)
	q.store.Find(key, func(res []byte, err error) {
		if err != nil {
			resChan <- ""
			errChan <- err
			return
		}
		if len(res) == 0 {
			innerChan := make(chan []byte)
			innerKey := successKey(proj, "")
			q.store.Find(innerKey, func(innerRes []byte, innerErr error) {
				if innerErr != nil {
					innerChan <- innerRes
					errChan <- err
					return
				}
				innerChan <- innerRes

			})
			res = <-innerChan
		}
		decRes := &lastSuccess{}
		errChan <- json.Unmarshal(res, decRes)
		resChan <- decRes.Commit
	})
	commit := <-resChan
	cErr := <-errChan
	return commit, cErr
}

// SaveLastSuccess saves the commits of the last successful build on a branch of a project
func (q *Queue) SaveLastSuccess(proj, branch, last string) error {
	errChan := make(chan error)
	key1 := successKey(proj, branch)
	key2 := successKey(proj, "")
	lastSucc := lastSuccess{Commit: last}
	q.store.Save(key1, lastSucc, func(err error) {
		errChan <- err
	})
	q.store.Save(key2, lastSucc, func(err error) {
		errChan <- err
	})
	for errCount := 0; errCount < 2; {
		err := <-errChan
		if err != nil {
			return err
		}
		errCount++
	}
	close(errChan)
	return nil
}

// Add adds a project to the queue
func (q *Queue) Add(proj, branch, prevCommit, currCommit string) {
	last, lastErr := q.GetLastSuccess(proj, branch)
	if lastErr != nil || len(last) == 0 {
		if lastErr != nil {
			fmt.Println("WARNING:", lastErr.Error())
		}
		last = prevCommit
	}
	newEntry := &QueueEntry{
		Project:    proj,
		Branch:     branch,
		PrevCommit: last,
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
				key := aliveKey{
					Project: next.Project,
					Branch:  next.Branch,
				}
				r.Alive[key] = next
				q.Queue = q.Queue[1:]
				return next
			}
		}
	}
	return nil
}

// SnapShot saves the queue's current state
func (q *Queue) SnapShot() error {
	errChan := make(chan error)
	q.store.Save("queue", q.Queue, func(err error) {
		errChan <- err
	})
	saveErr := <-errChan
	return saveErr
}
