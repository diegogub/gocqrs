package gocqrs

import (
	"errors"
	"github.com/diegogub/lib"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	UserEntity = "users"
)

type User struct {
	Username string    `json:"username"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
	Created  time.Time `json:"created`

	ActiveToken bool   `json:"activeToken"`
	Token       string `json:"token"`

	Groups []string               `json:"groups"`
	Data   map[string]interface{} `json:"data"`
}

type UserEventHandler struct {
}

const UserCreatedEvent = "UserCreated"
const UserTokenUpdated = "UserTokenUpdated"

func (uh UserEventHandler) EventName() []string {
	return []string{
		UserCreatedEvent,
		UserTokenUpdated,
	}
}

func (uh UserEventHandler) Handle(id, userid, role string, event Eventer, entity *Entity, replay bool) (StoreOptions, error) {
	var opt StoreOptions
	var err error
	switch event.GetType() {
	case UserCreatedEvent:
		var u User
		DecodeEvent(event, &u)
		u.Username = id

		// create new user
		opt.Create = true
		if !replay {
			err = u.Encrypt()
			if err != nil {
				return opt, err
			}
			u.Created = time.Now().UTC()
			u.GenerateToken()
			err = u.Valid()
			if err != nil {
				return opt, err
			}
			event.SetData("token", u.Token)
			event.SetData("password", u.Password)
		}

		entity.Data = ToMap(u)

	case UserTokenUpdated:
		if !replay {
			var u User
			u.GenerateToken()
			token := u.Token
			event.ClearData()
			event.SetData("token", token)
		} else {
			data := event.GetData()
			entity.Data["token"] = data["token"]
		}
	}

	return opt, err
}

func (uh UserEventHandler) CheckBase(e Eventer) bool {
	switch e.GetType() {
	case UserCreatedEvent:
		return true
	}
	return false
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

func (u *User) GenerateToken() {
	t := lib.NewLongId("t") + lib.NewLongId("d") + lib.NewLongId("4")
	u.Token = t
}
