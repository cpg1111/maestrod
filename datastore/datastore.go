package datastore

// ResultCallback is a callback that is expecting both a resulting data interface and an error
type ResultCallback func([]byte, error)

// NoResultCallback is a callback just expecting an error
type NoResultCallback func(error)

// Datastore is an interface for possible datastores
type Datastore interface {
	Save(key string, data interface{}, callback NoResultCallback)
	Find(queryStr string, callback ResultCallback)
	Remove(queryStr string, callback NoResultCallback)
	Update(queryStr string, update interface{}, callback NoResultCallback)
	FindAndUpdate(queryStr string, update interface{}, callback ResultCallback)
}
