package dom

import (
	"errors"
	nap "github.com/jmcvetta/napping"
	"net/url"
	"regexp"
)

var (
	STREAM_NO_MATCH = errors.New("stream don't match")
)

type PubRepository interface {
	Store(pub Pub) error
	List() []Pub
	Delete(id string) error
}

type Pub struct {
	Id     string `json:"id"`
	Regex  string `json:"regex"`
	regExp *regexp.Regexp
	Url    string `json:"url"`
}

func (p *Pub) Validate() error {
	if p.Id == "" {
		return errors.New("Invalid pub id")
	}

	if p.Url == "" {
		return errors.New("Invalid url")
	}

	u, err := url.Parse(p.Url)
	if err != nil {
		return err
	}

	if !u.IsAbs() {
		return errors.New("Invalid url,must be absolute")
	}

	_, err = regexp.Compile(p.Regex)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pub) Publish(streamid string, e Event) (int, error) {
	var err error

	if p.regExp == nil && p.Regex != "" {
		reg, _ := regexp.Compile(p.Regex)
		p.regExp = reg
	}

	if p.Regex != "" {
		match := p.regExp.MatchString(streamid)
		if !match {
			return 0, STREAM_NO_MATCH
		}
	}

	e.StreamId = streamid
	res, err := nap.Post(p.Url, e, nil, nil)
	if err != nil {
		return 0, err
	}

	return res.Status(), err
}
