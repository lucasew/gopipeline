package compose

import "reflect"

func implementsErr(t reflect.Type) bool {
    if t == nil {
        return false
    }
    e := reflect.TypeOf(&ErrCompose).Elem()
    return t.Implements(e)
}
