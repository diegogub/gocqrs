package gocqrs

type Role struct {
	Name string `json:"name"`
	// events and read cmd that can be executed ,default ALL
	Allowed []string `json:"allowed"`

	NotAllowed []string `json:"noAllowed"`
}

func NewRole(r string) *Role {
	var role Role
	if r == "" {
		panic("invalid role name")
	}
	role.Name = r
	role.Allowed = make([]string, 0)
	role.NotAllowed = make([]string, 0)
	return &role
}

func (r *Role) Can(e string) bool {
	allowed := false

	if len(r.Allowed) == 0 {
		allowed = true
	} else {
		for _, ae := range r.Allowed {
			if ae == e {
				allowed = true
				break
			}
		}
	}

	for _, de := range r.NotAllowed {
		if e == de {
			allowed = false
			break
		}
	}

	return allowed
}

func (r *Role) Allow(cmd ...string) {
	for _, c := range cmd {
		r.Allowed = append(r.Allowed, c)
	}
}

func (r *Role) NotAllow(cmd ...string) {
	for _, c := range cmd {
		r.NotAllowed = append(r.NotAllowed, c)
	}
}
