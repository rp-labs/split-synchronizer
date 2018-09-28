package redis

import (
	"encoding/json"
	"sync"

	"github.com/splitio/split-synchronizer/log"
	"github.com/splitio/split-synchronizer/splitio/api"
	redis "gopkg.in/redis.v5"
)

var elMutex = &sync.Mutex{}

// EventStorageAdapter implements EventStorage interface
type EventStorageAdapter struct {
	*BaseStorageAdapter
}

// NewEventStorageAdapter returns an instance of EventStorageAdapter
func NewEventStorageAdapter(clientInstance *redis.ClusterClient, prefix string) *EventStorageAdapter {
	prefixAdapter := &prefixAdapter{prefix: prefix}
	adapter := &BaseStorageAdapter{prefixAdapter, clientInstance}
	client := EventStorageAdapter{BaseStorageAdapter: adapter}
	return &client
}

// PopN returns elements given by LRANGE 0 items and perform a LTRIM items -1
func (r EventStorageAdapter) PopN(n int64) ([]api.RedisStoredEventDTO, error) {

	toReturn := make([]api.RedisStoredEventDTO, 0)

	elMutex.Lock()
	lrange := r.client.LRange(r.eventsListNamespace(), 0, n-1)
	if lrange.Err() != nil {
		log.Error.Println("Fetching events", lrange.Err().Error())
		elMutex.Unlock()
		return nil, lrange.Err()
	}
	totalFetchedEvents := int64(len(lrange.Val()))

	idxFrom := n
	if totalFetchedEvents < n {
		idxFrom = totalFetchedEvents
	}

	res := r.client.LTrim(r.eventsListNamespace(), idxFrom, -1)
	if res.Err() != nil {
		log.Error.Println("Trim events", res.Err().Error())
		elMutex.Unlock()
		return nil, res.Err()
	}
	elMutex.Unlock()

	//JSON unmarshal
	listOfEvents := lrange.Val()
	for _, se := range listOfEvents {
		storedEventDTO := api.RedisStoredEventDTO{}
		err := json.Unmarshal([]byte(se), &storedEventDTO)
		if err != nil {
			log.Error.Println("Error decoding event JSON", err.Error())
			continue
		}
		toReturn = append(toReturn, storedEventDTO)
	}

	return toReturn, nil
}

// Size return the value of LLEN
func (r EventStorageAdapter) Size() int64 {
	elMutex.Lock()
	llen := r.client.LLen(r.eventsListNamespace()) //LRange(r.eventsListNamespace(), 0, n-1)
	elMutex.Unlock()

	if llen.Err() != nil {
		log.Error.Println(llen.Err())
		return 0
	}

	return llen.Val()
}
