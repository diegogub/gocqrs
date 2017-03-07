package gocqrs

import (
	"errors"
	"github.com/diegogub/lib"
	"gopkg.in/asaskevich/govalidator.v4"
	"log"
)

type Validator interface {
	GetName() string
	Validate(e Entity) error
}

type SimpleValidator struct {
	validators []*EntityProperty
}

func (sv SimpleValidator) GetName() string {
	return "simple-validator"
}

func (sv SimpleValidator) Validate(e Entity) error {
	var err error
	for _, v := range sv.validators {
		for _, ep := range v.Validations {
			err = ep(e.Data[v.Name])
			if err != nil {
				return err
			}
		}
	}
	return err
}

func NewSimpleValidator(ep ...*EntityProperty) *SimpleValidator {
	var sv SimpleValidator
	sv.validators = make([]*EntityProperty, 0)

	for _, e := range ep {
		sv.validators = append(sv.validators, e)
	}

	return &sv
}

type EntityProperty struct {
	Name        string
	Validations []func(interface{}) error
}

func NewProperty(name string) *EntityProperty {
	var ep EntityProperty
	if name == "" {
		log.Fatal("invalid property")
	}

	ep.Name = name
	ep.Validations = make([]func(interface{}) error, 0)
	return &ep
}

func (ep *EntityProperty) String() *EntityProperty {
	valid := func(i interface{}) error {
		switch i.(type) {
		case string:
		default:
			return errors.New(ep.Name + " invalid string")
		}
		return nil
	}

	ep.Validations = append(ep.Validations, valid)
	return ep
}

func (ep *EntityProperty) ID() *EntityProperty {
	valid := func(i interface{}) error {
		switch i.(type) {
		case string:
			err := lib.ValidID(i.(string), true)
			if err != nil {
				return errors.New(ep.Name + " invalid id")
			}
		default:
			return errors.New(ep.Name + " invalid id")
		}
		return nil
	}

	ep.Validations = append(ep.Validations, valid)
	return ep
}

func (ep *EntityProperty) Email() *EntityProperty {
	valid := func(i interface{}) error {
		switch i.(type) {
		case string:
			isEmail := govalidator.IsEmail(i.(string))
			if !isEmail {
				return errors.New(ep.Name + " invalid email")
			}
		default:
			return errors.New(ep.Name + " invalid email")
		}
		return nil
	}

	ep.Validations = append(ep.Validations, valid)
	return ep
}

func (ep *EntityProperty) IP() *EntityProperty {
	valid := func(i interface{}) error {
		switch i.(type) {
		case string:
			isIP := govalidator.IsIP(i.(string))
			if !isIP {
				return errors.New(ep.Name + " invalid IP")
			}
		default:
			return errors.New(ep.Name + " invalid IP")
		}
		return nil
	}

	ep.Validations = append(ep.Validations, valid)
	return ep
}

func (ep *EntityProperty) URL() *EntityProperty {
	valid := func(i interface{}) error {
		switch i.(type) {
		case string:
			isURL := govalidator.IsURL(i.(string))
			if !isURL {
				return errors.New(ep.Name + " invalid url")
			}
		default:
			return errors.New(ep.Name + " invalid url")
		}
		return nil
	}

	ep.Validations = append(ep.Validations, valid)
	return ep
}

func (ep *EntityProperty) NotNull() *EntityProperty {
	valid := func(i interface{}) error {
		switch v := i.(type) {
		case string, []byte:
			if len(v.(string)) == 0 {
				return errors.New(ep.Name + " invalid value")
			}

		default:
			return errors.New(ep.Name + " invalid string")
		}
		return nil
	}

	ep.Validations = append(ep.Validations, valid)
	return ep
}
