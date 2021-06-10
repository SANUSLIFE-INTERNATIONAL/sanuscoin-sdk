// Copyright © 2021 The Sanuscoin Team

package context

import (
	ctx "context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"sanus/sanus-sdk/config"
)

type (
	// Context describes application's context interface.
	Context interface {
		ctx.Context
		Cancel()
		WgAdd(delta int)
		WgDone()
		WgWait()
		WithCancel() Context
		WithCancelWait() Context
		WithTimeout() Context
	}

	// context implements application's context interface.
	context struct {
		cancel    ctx.CancelFunc
		config    *config.Config
		parent    ctx.Context
		waitGroup *sync.WaitGroup
	}
)

var (
	// Make sure context implements context interface.
	_ Context = (*context)(nil)

	// Canceled is the error returned by Context.Err
	// when the context is canceled.
	Canceled = ctx.Canceled

	// DeadlineExceeded is the error returned by Context.Err
	// when the context's deadline passes.
	DeadlineExceeded = ctx.DeadlineExceeded
)

// NewContext creates context of the application.
func NewContext(cfg *config.Config) Context {
	cc, cancel := ctx.WithCancel(ctx.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, // to stop channel
			syscall.SIGINT,  // interrupt
			syscall.SIGQUIT, // quit
			syscall.SIGABRT, // aborted
			syscall.SIGKILL, // killed
			syscall.SIGTERM, // terminated
		)
		<-stop
		cancel()
	}()

	return &context{
		cancel:    cancel,
		config:    cfg,
		parent:    cc,
		waitGroup: new(sync.WaitGroup),
	}
}

func (c *context) Cancel() {
	c.cancel()
}

func (c *context) Deadline() (deadline time.Time, ok bool) {
	return c.parent.Deadline()
}

func (c *context) Done() <-chan struct{} {
	return c.parent.Done()
}

func (c *context) Err() error {
	return c.parent.Err()
}

func (c *context) Value(key interface{}) interface{} {
	return c.parent.Value(key)
}

// WgAdd adds delta to context's wait group.
func (c *context) WgAdd(delta int) {
	c.waitGroup.Add(delta)
}

// WgDone decrements context's wait group counter by one.
func (c *context) WgDone() {
	c.waitGroup.Done()
}

// WgWait blocks until resolved context wait group counter is zero.
func (c *context) WgWait() {
	c.waitGroup.Wait()
}

// WithCancel returns a copy of application context with a new done channel.
// The returned context's Done channel is closed when the returned cancel
// function is called or when the parent context's Done channel is closed,
// whichever happens first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func (c *context) WithCancel() Context {
	cc, cancel := ctx.WithCancel(c)
	return &context{
		cancel:    cancel,
		config:    c.config,
		parent:    cc,
		waitGroup: c.waitGroup,
	}
}

// WithCancelWait returns a copy of application context with a new done channel.
// The returned context's Done channel is closed when the returned cancel
// function is called or when the parent context's Done channel is closed,
// whichever happens first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func (c *context) WithCancelWait() Context {
	c.waitGroup.Add(1)
	return c.WithCancel()
}

// WithTimeout returns a copy of application context with timeout duration
// specified with application environment into the network configuration.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func (c *context) WithTimeout() Context {
	cc, cancel := ctx.WithCancel(c)
	return &context{
		cancel:    cancel,
		config:    c.config,
		parent:    cc,
		waitGroup: c.waitGroup,
	}
}
