package usecases

import (
	"bitbucket.org/dgub/evento/dom"
	"errors"
	"regexp"
)

func StreamVersion(streamid string) (uint64, error) {
	var version uint64
	var err error

	version, err = man.StreamRepo.Version(streamid)
	return version, err
}

func StreamExist(streamid string) bool {
	_, err := man.StreamRepo.Version(streamid)
	if err != nil {
		return false
	}
	return true
}

func StreamMatch(regex string) ([]dom.Stream, error) {
	list := make([]dom.Stream, 0)

	chStream := man.StreamRepo.GetAll()
	validStream, err := regexp.Compile(regex)
	if err != nil {
		return list, err
	}

	for s := range chStream {
		if validStream.MatchString(s.Id) {
			list = append(list, s)
		}
	}

	return list, nil

}

func (edb *EventoDB) DeleteStream(id string) error {
	edb.lock.Lock()
	defer edb.lock.Unlock()

	stream, err := man.StreamRepo.Get(id)
	if err == nil {
		if stream.Deleted != dom.StreamDeleted {
			stream.Deleted = dom.StreamDeleted
			man.StreamRepo.Save(*stream, true)
		} else {
			return errors.New("stream already deleted")
		}
	} else {
		return errors.New("stream not found")
	}

	return nil
}

func (edb *EventoDB) PurgeStream(id string) error {
	edb.lock.Lock()
	defer edb.lock.Unlock()

	return man.StreamRepo.Purge(id)
}

func GetPurges() []dom.PurgeStatus {
	list := make([]dom.PurgeStatus, 0)
	ch := man.StreamRepo.GetPurges()

	for ps := range ch {
		list = append(list, ps)
	}

	return list
}

func applyPurge(ps dom.PurgeStatus) error {
	return nil
}
