package fn

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func assertNotError(t *testing.T, err error) bool {
    if err != nil {
        t.Errorf("test returned error: %s", err)
        return true
    }
    return false
}

func assertEqual(t *testing.T, x interface {}, y interface{}) bool {
    if !reflect.DeepEqual(x, y) {
        t.Errorf("expected '%+v' got '%+v'", x, y)
        return true
    }
    return false
}

func spewValue(v ...reflect.Value) {
    iface := make([]interface{}, len(v))
    for i := 0; i < len(v); i++ {
        iface[i] = v[i].Interface()
    }
    spew.Dump(iface)
}
