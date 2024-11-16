package libraryErrors

import "fmt"

func Panicer(err any) {
	checker(err, func(a any) { panic(a) })
}

func Errorer(err error) {
	checker(err, func(a any) { libraryLogging.Error(a.(error).Error()) })
}

func checker(a any, f func(any)) {
	if a != nil {
		f(a)
	}
}

func WrapError(msg string, err error) error {
	if err != nil {
		return fmt.Errorf("%s -> %w", msg, err)
	}
	return nil
}
