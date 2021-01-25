package compose

import (
	"reflect"
	"testing"
)

func TestImplementsErr(t *testing.T) {
    assertEqual(t, implementsErr(reflect.TypeOf(ErrCompose)), true)
    assertEqual(t, implementsErr(reflect.TypeOf(nil)), false)
    assertEqual(t, implementsErr(reflect.TypeOf("error")), false)
}
