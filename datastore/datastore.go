package datastore

type ResultCallback func(interface{}, error)

type NoResultCallback func(error)

type Datastore interface {
	Save(key string, data interface{}, callback ResultCallback)
	Find(queryStr string, callback ResultCallback)
	Remove(queryStr string, callback NoResultCallback)
	Update(queryStr string, update interface{}, callback NoResultCallback)
	FindAndUpdate(queryStr string, update interface{}, callback ResultCallback)
}
