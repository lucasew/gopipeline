package fn

import "reflect"

func ImplementsErr(t reflect.Type) bool {
    if t == nil {
        return false
    }
    e := reflect.TypeOf(&ErrCompose).Elem()
    return t.Implements(e)
}
