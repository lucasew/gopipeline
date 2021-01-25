package compose

import (
	"fmt"
	"reflect"
)

var (
    // ErrCompose Error that happen while Compose typechecks the stages and the function pointer
    ErrCompose = fmt.Errorf("compose: ")
)

// Compose build a function out that is the composition of the stages function like a pipeline, out must be a pointer to a empty function created using the var declaration, errors do early returns and are not passed to the next stage. The previous stage must return the input of the current stage. Any inconsistency will return a error.
func Compose(out interface{}, stages ...interface{}) (err error) {
    // check basic typing constraints for out
    if out == nil {
        return fmt.Errorf("%w you passed nil as out parameter", ErrCompose)
    }
    // shortcut variables
    vOut := reflect.ValueOf(out)
    tOut := vOut.Type()
    if vOut.Kind() != reflect.Ptr || vOut.Elem().Kind() != reflect.Func {
        return fmt.Errorf("%w you must pass a pointer to a function for out parameter", ErrCompose)
    }
    // if no stages specified then problem
    if len(stages) == 0 {
        return fmt.Errorf("%w no stages specified", ErrCompose)
    }
    // if only one stage check typing and associate it to out
    if len(stages) == 1 {
        if reflect.TypeOf(stages[0]).AssignableTo(vOut.Elem().Type()) {
            vOut.Elem().Set(reflect.ValueOf(stages[0]))
            return nil
        }
        return fmt.Errorf("%w cant assign first stage to out in monostage composition", ErrCompose)
    }
    // check if all stages are functions and if some of them return errors
    haveError := false
    for nth, stage := range stages { //checkup
        vStage := reflect.ValueOf(stage)
        if vStage.Kind() != reflect.Func {
            return fmt.Errorf("%w error near stage %d: not a function but a %s", ErrCompose, nth + 1, vStage.Type().String())
        }
        outs := vStage.Type().NumOut()
        if outs == 0 {
            continue
        }
        if implementsErr(vStage.Type().Out(outs - 1)) {
            haveError = true
        }
    }
    // if the functions return errors then check if the out function can return that error
    if haveError && (tOut.Elem().NumOut() == 0 || !implementsErr(tOut.Elem().Out(tOut.Elem().NumOut() - 1))) {
        return fmt.Errorf("%w provided target function must return at least a error", ErrCompose)
    }
    // validate stage typings
    input := make([]reflect.Type, 0, tOut.Elem().NumIn())
    for i := 0; i < tOut.Elem().NumIn(); i++ {
        input = append(input, tOut.Elem().In(i))
    }
    for nth, stage := range stages {
        // spew.Dump(input)
        tStage := reflect.TypeOf(stage)
        if tStage.NumIn() != len(input) {
            return fmt.Errorf("%w %dth stage expects %d arguments but the previous stage send %d", ErrCompose, nth + 1, tStage.NumIn(), len(input))
        }
        for i := 0; i < tStage.NumIn(); i++ {
            if !tStage.In(i).AssignableTo(input[i]) {
                return fmt.Errorf("%w %dth stage expects as its %d parameter a %s but the provided %s cant be assigned to the expected type", ErrCompose, nth + 1, i + 1, tStage.In(i).String(), input[i].String())
            }
        }
        outs := tStage.NumOut()
        if implementsErr(tStage.Out(outs - 1)) {
            outs--
        }
        input = make([]reflect.Type, 0, outs)
        for i := 0; i < outs; i++ {
            input = append(input, tStage.Out(i))
        }
    }
    for i := 0; i < len(input); i++ {
        if !input[i].AssignableTo(tOut.Elem().Out(i)) {
            return fmt.Errorf("%w %dth parameter of the last stage of the pipeline (to return the data), that is a %s cant be assigned to %s", ErrCompose, i + 1, input[i].String(), tOut.Elem().Out(i).String())
        }
    }
    // now input is the return of the last stage, that is the return of the composed function
    earlyRet := func (err error) []reflect.Value {
        outSize := tOut.Elem().NumOut()
        ret := make([]reflect.Value, 0, outSize)
        for i := 0; i < outSize; i++ {
            ret = append(ret, reflect.New(tOut.Elem().Out(i)).Elem())
        }
        if haveError {
            ret[outSize - 1] = reflect.ValueOf(err)
        }
        return ret
    }
    generatedFn := reflect.MakeFunc(tOut.Elem(), func (args []reflect.Value) ([]reflect.Value) {
        data := args
        for _, stage := range stages {
            vStage := reflect.ValueOf(stage)
            data = vStage.Call(data)
            lastRet := data[len(data) - 1]
            if implementsErr(lastRet.Type()) {
                if !lastRet.IsNil() {
                    er := earlyRet(lastRet.Interface().(error))
                    return er
                }
                data = data[:len(data) - 1]
            }
        }
        if haveError {
            var nilerr error
            data = append(data, reflect.ValueOf(&nilerr).Elem())
        }
        return data
    })
    vOut.Elem().Set(generatedFn)
    // at the end of the process check if the function pointer was altered
    if vOut.Elem().IsNil() {
        return fmt.Errorf("%w function pointer is still nil", ErrCompose)
    }
    return nil
}
