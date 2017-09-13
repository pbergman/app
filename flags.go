package app

import "reflect"

// FlagSet is a simple interface where a flag lib has to
// implement. Created this so you can use any flag lib
// that is based on the internal flab lib. (go-flags, pflags etc.)
type FlagSet interface {
	Parse([]string) error
	Args() []string
}

// setFlagsUsage is a simple wrapper to set the Usage property
// because interface can not hold properties we have to
// do it with the reflection lib
func setFlagsUsage(f FlagSet, u func()) {
	elem := reflect.ValueOf(f).Elem()
	if usage := elem.FieldByName("Usage"); usage.IsValid() {
		usage.Set(reflect.ValueOf(u))
	}
}
