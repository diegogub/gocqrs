package dom

import (
	"errors"
	"strings"
	"time"
)

const (
	SHARD_SEP     = "|"
	REF_SEP       = "-"
	StreamDeleted = "d"
	StreamPurge   = "p"
	Purged        = "e"
)

type StreamRepository interface {
	Create(s Stream) error
	// save stream into database
	Save(s Stream, sync bool) error
	// get stream
	Get(id string) (*Stream, error)
	// get all streams from database
	GetAll() chan Stream
	GetPurges() chan PurgeStatus
	// get current version
	Version(id string) (uint64, error)
	// mark as deleted
	Delete(s *Stream) error
	// Truncate events
	Purge(id string) error
	// Do real database purge
	DoPurge(ps PurgeStatus) error
}

type Stream struct {
	// StreamId Example:
	// (shardId)|(ref1)-(ref2-)-...-(refn)
	//    example: ny3|campaign_tracking-1231230912901
	Id      string    `json:"id"`
	Created time.Time `json:"c"`
	Version uint64    `json:"v"`
	Deleted string    `json:"d,omitempty"`
	// Purge status
	Purge  string `json:"p,omitempty"`
	Purged string `json:"p,omitempty"`
}

type PurgeStatus struct {
	Id      string `json:"id"`
	Current uint64 `json:"cv"`
	Total   uint64 `json:"t"`
}

func NewPurgeStatus(s Stream) *PurgeStatus {
	var ps PurgeStatus
	ps.Id = s.Id
	ps.Current = 0
	ps.Total = s.Version
	return &ps
}

func NewStream(id string, sync bool, postback []string, log bool) (*Stream, error) {
	var s Stream
	if id == "" || len(id) > 150 {
		return nil, errors.New("Invalid stream id")
	}

	s.Id = id
	s.Created = time.Now().UTC()
	return &s, nil
}

func BuildStreamId(shardId string, refs ...string) (string, error) {
	var streamId string

	if strings.Contains(shardId, "$") {
		return "", errors.New("invalid shardid")
	}

	if strings.Contains(shardId, "|") {
		return "", errors.New("invalid shardid")
	}

	for _, r := range refs {
		if strings.Contains(r, "-") {
			return "", errors.New("invalid shardid")
		}

		if strings.Contains(r, "$") {
			return "", errors.New("invalid shardid")
		}
	}

	ref := strings.Join(refs, REF_SEP)
	if shardId != "" {
		streamId = shardId + SHARD_SEP + ref
	} else {
		streamId += ref
	}
	return streamId, nil
}

func ParseStreamId(streamid string) (shardid string, refs string, err error) {
	if streamid == "" {
		return shardid, refs, errors.New("Invalid StreamId")
	}
	parts := strings.Split(streamid, "|")
	if len(parts) == 2 {
		shardid = parts[0]
		refs = parts[1]
	} else {
		refs = parts[1]
	}

	return shardid, refs, err
}

// Do postback to http endpoints
func (s *Stream) Postback(e Event) error {
	return nil
}
