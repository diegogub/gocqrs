package lib

import (
	"errors"
	"github.com/rs/xid"
	"regexp"
	"strings"
)

var idRegex *regexp.Regexp

var (
	InvalidEventIDError = errors.New("Invalid ID")
)

func init() {
	r, err := regexp.Compile(`^[A-Za-z0-9]+[\w\-\:\.@]*$`)
	if err != nil {
		panic(err)
	}
	idRegex = r
}

func BuildPrefix(ids ...string) string {
	return strings.Join(ids, "")
}

// Prefix should be unique within clusters
func NewLongId(prefix string) string {
	var id string
	//var base int64 = 10 ^ 6
	guid := xid.New()

	id = guid.String()
	/*
		//random number
		num := base + r.Int63n(999999999)
		// ticks
		nsec := uint64(time.Now().UnixNano() / 100)
		id = prefix + strconv.FormatUint(uint64(num), 32) + strconv.FormatUint(nsec, 32)
	*/

	return id
}

func NewShortId(prefix string) string {
	var id string
	id = xid.New().String()
	return id
}

func ValidID(id string, notNull bool) error {
	if notNull && id == "" {
		return InvalidEventIDError
	} else {
		if id == "" {
			return nil
		}
	}
	if len(id) > 99 {
		return InvalidEventIDError
	}

	if !idRegex.MatchString(id) {
		return InvalidEventIDError
	}

	return nil
}
