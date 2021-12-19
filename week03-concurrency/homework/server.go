package homework

import (
	"context"
	"os"
	"os/signal"
)

// Server creates function for canceling execution.
// First function passed in should block.
// Second function passed in should unblock first function.
func Server(serve func() error, shutdown func() error) RunFunc {
	return func(stop <-chan struct{}) error {
		done := make(chan error)
		defer close(done)

		go func() {
			done <- serve()
		}()

		select {
		case err := <-done:
			return err
		case <-stop:
			err := shutdown()
			if err == nil {
				err = <-done
			} else {
				<-done
			}
			return err
		}
	}
}

// Signal creates function for canceling execution using os signal.
func Signal(sig ...os.Signal) RunFunc {
	return func(stop <-chan struct{}) error {
		if len(sig) == 0 {
			sig = append(sig, os.Interrupt)
		}
		done := make(chan os.Signal, len(sig))
		defer close(done)

		signal.Notify(done, sig...)
		defer signal.Stop(done)

		select {
		case <-stop:
		case <-done:
		}
		return nil
	}
}

// RunFunc is a function to execute with other related functions in its own goroutine.
// The closure of the channel passed to RunFunc should trigger return.
type RunFunc func(<-chan struct{}) error

// Group is a group of related goroutines.
// The zero value for a Group is fully usable without initialization.
type Group struct {
	fns []RunFunc
}

// Add adds a function to the Group.
// The function will be exectuted in its own goroutine when Run is called.
// Add must be called before Run.
func (g *Group) Add(fn RunFunc) {
	g.fns = append(g.fns, fn)
}

// Run executes each function registered via Add in its own goroutine.
// Run blocks until all functions have returned, then returns the first non-nil error (if any) from them.
// The first function to return will trigger the closure of the channel passed to each function, which should in turn, return.
func (g *Group) Run() error {
	if len(g.fns) == 0 {
		return nil
	}

	stop := make(chan struct{})
	done := make(chan error, len(g.fns))
	defer close(done)

	for _, fn := range g.fns {
		go func(fn RunFunc) {
			done <- fn(stop)
		}(fn)
	}

	var err error
	for i := 0; i < cap(done); i++ {
		if err == nil {
			err = <-done
		} else {
			<-done
		}
		if i == 0 {
			close(stop)
		}
	}
	return err
}

// Context creates function for canceling execution using context.
func Context(ctx context.Context) RunFunc {
	return func(stop <-chan struct{}) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-stop:
			return nil
		}
	}
}
