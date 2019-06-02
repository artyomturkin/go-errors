package errors_test

import (
	"fmt"
	errors "github.com/artyomturkin/go-errors"
	"sync"
	"testing"
	"time"
)

///////////////////     Example     ///////////////////
func ExampleChannel() {
	errChan := &errors.Channel{}

	// sync goroutines
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// create subscription
	errCh := errChan.Errors()

	// print errors
	go func() {
		for err := range errCh {
			fmt.Printf("%v\n", err)
		}
		wg.Done()
	}()

	// publish an error
	errChan.Publish(fmt.Errorf("new error"))

	// close channel
	errChan.Close()

	// wait for goroutines to finish
	wg.Wait()

	// Output: new error
}

///////////////////     Tests     ///////////////////
func createErrors(e *errors.Channel, count int) {
	for index := 0; index < count; index++ {
		e.Publish(fmt.Errorf("%d", index))
	}
}

func consumeErrors(e *errors.Channel) (*[]error, func()) {
	errs := &[]error{}
	errsch := e.Errors()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for err := range errsch {
			*errs = append(*errs, err)
		}
		wg.Done()
	}()

	return errs, wg.Wait
}

func TestErrors_SebscribeBefore(t *testing.T) {
	count := 10
	e := &errors.Channel{}

	errs, wait := consumeErrors(e)
	createErrors(e, count)
	e.Close()
	wait()

	if len(*errs) != count {
		t.Errorf("not all errors propogated. Want %d, got %d", count, len(*errs))
	}
}

func TestErrors_SebscribeAfter(t *testing.T) {
	count := 10
	e := &errors.Channel{}

	createErrors(e, count)
	time.Sleep(time.Millisecond)
	errs, wait := consumeErrors(e)
	e.Close()
	wait()

	if len(*errs) != 0 {
		t.Errorf("not all errors propogated. Want %d, got %d", count, len(*errs))
	}
}

func TestErrors_SebscribeAfterClose(t *testing.T) {
	count := 10
	e := &errors.Channel{}

	createErrors(e, count)
	e.Close()
	time.Sleep(time.Millisecond)
	errs, wait := consumeErrors(e)
	wait()

	if len(*errs) != 0 {
		t.Errorf("not all errors propogated. Want %d, got %d", count, len(*errs))
	}
}
