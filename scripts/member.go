package scripts

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Member field of a host value
type Member struct {
	host   Token
	member string
}

// Execute returns member value of host
func (member *Member) Execute(variables *Variables) (interface{}, error) {
	hostvalue, err := member.host.Execute(variables)
	if err != nil {
		return nil, err
	}

	if hostvalue == nil {
		return nil, errors.New("Null reference")
	}

	membername := strings.ToLower(member.member)
	typevalue := reflect.ValueOf(hostvalue).Elem()
	hosttype := reflect.TypeOf(hostvalue).Elem()
	for i := 0; i < hosttype.NumField(); i++ {
		field := hosttype.Field(i)
		if strings.ToLower(field.Name) == membername {
			return typevalue.Field(i).Interface(), nil
		}
	}

	return nil, fmt.Errorf("Member with name '%s' not found on '%v'", member.member, hostvalue)
}
