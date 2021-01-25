package compose

import (
	"reflect"
	"testing"
)

func TestImplementsErr(t *testing.T) {
    assertEqual(t, ImplementsErr(reflect.TypeOf(ErrCompose)), true)
    assertEqual(t, ImplementsErr(reflect.TypeOf(nil)), false)
    assertEqual(t, ImplementsErr(reflect.TypeOf("error")), false)
}
