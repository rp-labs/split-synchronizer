package wrappers

import (
	"errors"

	"github.com/splitio/split-synchronizer/log"
	"github.com/splitio/split-synchronizer/splitio/storage/boltdb"
	"github.com/splitio/split-synchronizer/splitio/storage/boltdb/collections"
	"github.com/splitio/split-synchronizer/splitio/storageDTOs"
)

var errSegmentStorageNotImplementedMethod = errors.New("Method has not been implemented yet")

// SegmentChangesWrapper implements SegmentStorage interface
type SegmentChangesWrapper struct {
	segmentCollection collections.SegmentChangesCollection
}

// NewSegmentChangesWrapper returns a new instance of SegmentChangesWrapper
func NewSegmentChangesWrapper() *SegmentChangesWrapper {
	return &SegmentChangesWrapper{segmentCollection: collections.NewSegmentChangesCollection(boltdb.DBB)}
}

// RegisteredSegmentNames returns a list of segment names
func (s *SegmentChangesWrapper) RegisteredSegmentNames() ([]string, error) {
	segments, err := s.segmentCollection.FetchAll()
	if err != nil {
		log.Error.Println("Error fetching segments from boldb", err)
		return nil, err
	}

	toReturn := make([]string, 0)
	for _, segment := range segments {
		toReturn = append(toReturn, segment.Name)
	}

	return toReturn, nil
}

// AddToSegment not implemented due this wrapper is only for dashboard
func (s *SegmentChangesWrapper) AddToSegment(segmentName string, keys []string) error {
	return errSegmentStorageNotImplementedMethod
}

// RemoveFromSegment not implemented due this wrapper is only for dashboard
func (s *SegmentChangesWrapper) RemoveFromSegment(segmentName string, keys []string) error {
	return errSegmentStorageNotImplementedMethod
}

// SetChangeNumber not implemented due this wrapper is only for dashboard
func (s *SegmentChangesWrapper) SetChangeNumber(segmentName string, changeNumber int64) error {
	return errSegmentStorageNotImplementedMethod
}

// ChangeNumber returns the change number of the segment
func (s *SegmentChangesWrapper) ChangeNumber(segmentName string) (int64, error) {
	segment, err := s.segmentCollection.Fetch(segmentName)
	if err != nil {
		log.Error.Println("Error fetching data for segment", segmentName)
		return 0, err
	}

	changeNumber := int64(0)
	for _, k := range segment.Keys {
		if k.ChangeNumber > changeNumber {
			changeNumber = k.ChangeNumber
		}
	}

	return changeNumber, nil
}

// CountActiveKeys return the active keys number
func (s *SegmentChangesWrapper) CountActiveKeys(segmentName string) (int64, error) {
	segment, err := s.segmentCollection.Fetch(segmentName)
	if err != nil {
		log.Error.Println("Error fetching data for segment", segmentName)
		return 0, err
	}

	changeNumber := int64(0)
	removedKeys := int64(0)
	addedKeys := int64(0)
	for _, k := range segment.Keys {
		if k.ChangeNumber > changeNumber {
			changeNumber = k.ChangeNumber
		}

		if k.Removed {
			removedKeys++
		} else {
			addedKeys++
		}
	}

	return addedKeys, nil
}

// Keys returns a list of keys
func (s *SegmentChangesWrapper) Keys(segmentName string) ([]storageDTOs.SegmentKeyDTO, error) {
	segment, err := s.segmentCollection.Fetch(segmentName)
	if err != nil {
		log.Error.Println("Error fetching data for segment", segmentName)
		return nil, err
	}

	toReturn := make([]storageDTOs.SegmentKeyDTO, 0)

	for _, k := range segment.Keys {
		toReturn = append(toReturn, storageDTOs.SegmentKeyDTO{
			Name:         k.Name,
			LastModified: k.ChangeNumber,
			Removed:      k.Removed,
		})
	}

	return toReturn, nil
}

// CountRemovedKeys return the removed keys number
func (s *SegmentChangesWrapper) CountRemovedKeys(segmentName string) (int64, error) {
	segment, err := s.segmentCollection.Fetch(segmentName)
	if err != nil {
		log.Error.Println("Error fetching data for segment", segmentName)
		return 0, err
	}

	changeNumber := int64(0)
	removedKeys := int64(0)
	addedKeys := int64(0)
	for _, k := range segment.Keys {
		if k.ChangeNumber > changeNumber {
			changeNumber = k.ChangeNumber
		}

		if k.Removed {
			removedKeys++
		} else {
			addedKeys++
		}
	}

	return removedKeys, nil
}
