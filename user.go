package gocqrs

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Username string    `json:"username"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
	Created  time.Time `json:"created`

	ActiveToken bool   `json:"activeToken"`
	Token       string `json:"token"`

	Data map[string]interface{} `json:"data"`
}

type UserEventHandler struct {
}

const UserCreatedEvent = "UserCreated"

func (uh UserEventHandler) EventName() []string {
	return []string{
		UserCreatedEvent,
	}
}

func (uh UserEventHandler) Handle(id string, event Eventer, entity *Entity) (StoreOptions, error) {
	var opt StoreOptions
	var err error
	switch event.GetType() {
	case UserCreatedEvent:
		var u User
		DecodeEvent(event, &u)
		u.Username = id

		// create new stream
		opt.Create = true
		err = u.Encrypt()
		if err != nil {
			return opt, err
		}
		u.Created = time.Now().UTC()

		err = u.Valid()
		if err != nil {
			return opt, err
		}

		event.SetData("password", u.Password)
		entity.Data = ToMap(u)
	}

	return opt, err
}

func (u *User) Valid() error {
	var err error
	if u.Username == "" {
		return errors.New("Invalid username")
	}

	_, ok := runningApp.Roles[u.Role]
	if !ok {
		return errors.New("Invalid role")
	}

	return err
}

func (u *User) Encrypt() error {
	if u.Password == "" {
		return errors.New("Invalid password")
	}
	b, err := bcrypt.GenerateFromPassword([]byte(u.Password), 11)
	u.Password = string(b)
	return err
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) CheckToken(token string) error {
	if u.Token == token {
		return nil
	}
	return errors.New("Invalid token")
}
