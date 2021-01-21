package fn

import (
	"fmt"
	"reflect"
)

var (
    ErrCompose = fmt.Errorf("compose:")
    nilerr error = nil
)

func Compose(out interface{}, stages ...interface{}) (err error) {
    // check basic typing constraints for out
    if out == nil {
        return fmt.Errorf("%w you passed nil as out parameter", ErrCompose)
    }
    // shortcut variables
    v_out := reflect.ValueOf(out)
    t_out := v_out.Type()
    if v_out.Kind() != reflect.Ptr || v_out.Elem().Kind() != reflect.Func {
        return fmt.Errorf("%w you must pass a pointer to a function for out parameter", ErrCompose)
    }
    // if no stages specified then problem
    if len(stages) == 0 {
        return fmt.Errorf("%w no stages specified", ErrCompose)
    }
    // if only one stage check typing and associate it to out
    if len(stages) == 1 {
        if reflect.TypeOf(stages[0]).AssignableTo(v_out.Elem().Type()) {
            v_out.Elem().Set(reflect.ValueOf(stages[0]))
            return nil
        } else {
            return fmt.Errorf("%w cant assign first stage to out in monostage composition", ErrCompose)
        }
    }
    // check if all stages are functions and if some of them return errors
    haveError := false
    for nth, stage := range stages { //checkup
        v_stage := reflect.ValueOf(stage)
        if v_stage.Kind() != reflect.Func {
            return fmt.Errorf("%w error near stage %d: not a function but a %s", ErrCompose, nth + 1, v_stage.Type().String())
        }
        outs := v_stage.Type().NumOut()
        if outs == 0 {
            continue
        }
        if ImplementsErr(v_stage.Type().Out(outs - 1)) {
            haveError = true
        }
    }
    // if the functions return errors then check if the out function can return that error
    if haveError && (t_out.Elem().NumOut() == 0 || !ImplementsErr(t_out.Elem().Out(t_out.Elem().NumOut() - 1))) {
        return fmt.Errorf("%w provided target function must return at least a error", ErrCompose)
    }
    // validate stage typings
    input := make([]reflect.Type, 0, t_out.Elem().NumIn())
    for i := 0; i < t_out.Elem().NumIn(); i++ {
        input = append(input, t_out.Elem().In(i))
    }
    for nth, stage := range stages {
        // spew.Dump(input)
        t_stage := reflect.TypeOf(stage)
        if t_stage.NumIn() != len(input) {
            return fmt.Errorf("%w %dth stage expects %d arguments but the previous stage send %d", ErrCompose, nth + 1, t_stage.NumIn(), len(input))
        }
        for i := 0; i < t_stage.NumIn(); i++ {
            if !t_stage.In(i).AssignableTo(input[i]) {
                return fmt.Errorf("%w %dth stage expects as its %d parameter a %s but the provided %s cant be assigned to the expected type", ErrCompose, nth + 1, i + 1, t_stage.In(i).String(), input[i].String())
            }
        }
        outs := t_stage.NumOut()
        if ImplementsErr(t_stage.Out(outs - 1)) {
            outs -= 1
        }
        input = make([]reflect.Type, 0, outs)
        for i := 0; i < outs; i++ {
            input = append(input, t_stage.Out(i))
        }
    }
    for i := 0; i < len(input); i++ {
        if !input[i].AssignableTo(t_out.Elem().Out(i)) {
            return fmt.Errorf("%w %dth parameter of the last stage of the pipeline (to return the data), that is a %s cant be assigned to %s", ErrCompose, i + 1, input[i].String(), t_out.Elem().Out(i).String())
        }
    }
    // now input is the return of the last stage, that is the return of the composed function
    earlyRet := func (err error) []reflect.Value {
        out_size := t_out.Elem().NumOut()
        ret := make([]reflect.Value, 0, out_size)
        for i := 0; i < out_size; i++ {
            ret = append(ret, reflect.New(t_out.Elem().Out(i)).Elem())
        }
        if haveError {
            ret[out_size - 1] = reflect.ValueOf(err)
        }
        return ret
    }
    generatedFn := reflect.MakeFunc(t_out.Elem(), func (args []reflect.Value) ([]reflect.Value) {
        data := args
        for _, stage := range stages {
            v_stage := reflect.ValueOf(stage)
            data = v_stage.Call(data)
            lastRet := data[len(data) - 1]
            if ImplementsErr(lastRet.Type()) {
                if !lastRet.IsNil() {
                    er := earlyRet(lastRet.Interface().(error))
                    return er
                }
                data = data[:len(data) - 1]
            }
        }
        if haveError {
            data = append(data, reflect.ValueOf(&nilerr).Elem())
        }
        return data
    })
    v_out.Elem().Set(generatedFn)
    // at the end of the process check if the function pointer was altered
    if v_out.Elem().IsNil() {
        return fmt.Errorf("%w function pointer is still nil", ErrCompose)
    }
    return nil
}
