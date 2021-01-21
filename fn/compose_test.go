package fn

import (
	"fmt"
	"testing"
)


func TestComposeIdentity(t *testing.T) {
    sum2 := func(x int) int {
        return x + 2
    }
    var composed func(int) int
    if assertNotError(t, Compose(&composed, sum2)) {return}
    if assertEqual(t, composed(2), 4) {return}
}

func TestComposeFn(t *testing.T) {
    sum2 := func(a int) int {
        return a + 2
    }
    var composed func(int)int
    if assertNotError(t, Compose(&composed, sum2, sum2)) {return}
    if assertEqual(t, composed(2), 6) {return}
}

func TestComposeFnError(t *testing.T) {
    fa := func(a int) (int, error) {
        return a + 2, nil
    }
    fb := func(b int) int {
        return b + 2
    }
    var composed func(int)(int, error)
    if assertNotError(t, Compose(&composed, fa, fb)) {return}
    res, err := composed(2)
    if assertNotError(t, err) {return}
    if assertEqual(t, res, 6) {return}
}


func TestComposeFnThird(t *testing.T) {
    sum2 := func(a int) int {
        return a + 2
    }
    sum3 := func(a int) int {
        return a + 3
    }
    var composed func(int)int
    if assertNotError(t, Compose(&composed, sum2, sum2, sum3)) {return}
    if assertEqual(t, composed(2), 9) {return}
}

func TestComposeFnMultipleParameters(t *testing.T) {
    sum := func (x, y int) int {
        return x + y
    }
    sum2 := func (x int) int {
        return x + 2
    }
    var composed func(int, int) int
    if assertNotError(t, Compose(&composed, sum, sum2)) {return}
    if assertEqual(t, composed(2,2), 6) {return}
}

func TestComposeFnMultipleParametersErr(t *testing.T) {
    sum := func (x, y int) (int, error) {
        return x + y, nil
    }
    sum2 := func (x int) int {
        return x + 2
    }
    var composed func(int, int) (int, error)
    if assertNotError(t, Compose(&composed, sum, sum2)) {return}
    r, err := composed(2, 2)
    if assertNotError(t, err) {return}
    if assertEqual(t, r, 6) {return}
}

func TestComposeFnMultipleReturns(t *testing.T) {
    ret := func(x int) (int, int) {
        return x*2, x*4
    }
    sum := func(x, y int) int {
        return x + y
    }
    var composed func(int) int
    if assertNotError(t, Compose(&composed, ret, sum)) {return}
    if assertEqual(t, composed(2), 12) {return}
}

func TestComposeFnMultipleReturnsAndError(t *testing.T) {
    ret := func(x int) (int, int, error) {
        return x*2, x*4, nil
    }
    sum := func(x, y int) int {
        return x + y
    }
    var composed func(int) (int, error)
    if assertNotError(t, Compose(&composed, ret, sum)) {return}
    r, err := composed(2)
    if assertNotError(t, err) {return}
    if assertEqual(t, r, 12) {return}
}

func TestComposeFnMultipleReturnsAndErrorNotNil(t *testing.T) {
    ret := func(x int) (int, int, error) {
        return x*2, x*4, fmt.Errorf("no problem")
    }
    sum := func(x, y int) int {
        return x + y
    }
    var composed func(int) (int, error)
    if assertNotError(t, Compose(&composed, ret, sum)) {return}
    r, err := composed(2)
    if err == nil {
        t.Errorf("composed should return error")
        return
    }
    if assertEqual(t, r, 0) {return} // dont setup the return value on error
}

func TestComposePassingNonsense(t *testing.T) {
    err := Compose(2, 2, 2)
    if err == nil {
        t.Errorf("nonsense was accepted")
    }
}

func TestComposePassingNonPointerFunction(t *testing.T) {
    sum2 := func(x int) int {
        return x + 2
    }
    var composed func(int) int
    err := Compose(composed, sum2)
    if err == nil {
        t.Errorf("non pointer function was accepted")
    }
}
