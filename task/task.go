package task

import "context"

type Task func(context.Context) error

func NewTask(t interface{}) Task {
    switch fn := t.(type) {
    case func():
        return func(context.Context) error {
            fn()
            return nil
        }
    case func() error:
        return func(context.Context) error {
            return fn()
        }
    case func(context.Context) error:
        return func(ctx context.Context) error {
            return fn(ctx)
        }
    case []func():
        return func(ctx context.Context) error {
            for _, f := range fn {
                select {
                case <-ctx.Done():
                    return context.Canceled
                default:
                    f()
                }
            }
            return nil
        }
    case []func()error:
        return func(ctx context.Context) error {
            for _, f := range fn {
                select {
                case <-ctx.Done():
                    return context.Canceled
                default:
                    err := f()
                    if err != nil {
                        return err
                    }
                }
            }
            return nil
        }
    case []func(context.Context)error:
        return func(ctx context.Context) error {
            for _, f := range fn {
                select {
                case <-ctx.Done():
                    return context.Canceled
                default:
                    err := f(ctx)
                    if err != nil {
                        return err
                    }
                }
            }
            return nil
        }
    default:
        panic("NewTask with something that can't be converted to a task")
    }
}



