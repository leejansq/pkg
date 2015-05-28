package enginer

import (
	"fmt"
	"io"
	//"io/ioutil"
	"os"
	"sync"
	"time"
)

type Handle func(*Job) error

type Enginer struct {
	id           string
	Handles      map[string]Handle
	stdout       io.Writer
	stdin        io.Reader
	stderr       io.Writer
	l            sync.Mutex
	tasks        sync.WaitGroup
	shutdownWait sync.WaitGroup
	logging      bool
	onShutdown   []func()
}

var globalhandle map[string]Handle

func GlobalRegister(name string, h Handle) error {
	if _, fined := globalhandle[name]; fined {
		return fmt.Errorf("The global command %s is exists! ", name)
	}
	globalhandle[name] = h
	return nil
}

func (r *Enginer) Register(name string, handle Handle) error {
	if _, fined := r.Handles[name]; fined {
		return fmt.Errorf("The command %s is exists! ", name)
	}
	r.Handles[name] = handle
	return nil
}

func (r *Enginer) Unregister(name string) {
	delete(r.Handles, name)
}

func New() *Enginer {
	eng := &Enginer{
		Handles: make(map[string]Handle),
		stdin:   os.Stdin,
		stdout:  os.Stdout,
		stderr:  os.Stderr,
		logging: true,
	}
	for name, h := range globalhandle {
		if err := eng.Register(name, h); err != nil {
			fmt.Errorf("global register err is %v", err)
			continue
		}
	}
	return eng
}

func (r *Enginer) Job(name string, args ...string) *Job {
	job := &Job{
		Eng:      r,
		Name:     name,
		Args:     args,
		Stdin:    NewInput(),
		Stdout:   NewOutput(),
		Stderr:   NewOutput(),
		closeIO:  true,
		Env:      make(map[string]interface{}),
		canclled: make(chan struct{}),
	}
	if r.logging {
		job.Stderr.Add(nopWriteCloser(r.stderr))
	}

	// Catchall is shadowed by specific Register.
	if handler, exists := r.Handles[name]; exists {
		job.handler = handler
	}
	return job
}

func (r *Enginer) Onshutdown(f func()) {
	r.l.Lock()
	r.onShutdown = append(r.onShutdown, f)
	r.shutdownWait.Add(1)
	r.l.Unlock()
}

func (r *Enginer) Shutdown() {
	r.l.Lock()
	taskdone := make(chan struct{})
	go func() {
		r.tasks.Wait()
		close(taskdone)
	}()
	select {
	case <-time.After(time.Second * 5):
	case <-taskdone:
	}
	for _, f := range r.onShutdown {
		go func(fc func()) {
			fc()
			r.shutdownWait.Done()
		}(f)
	}

	shutdone := make(chan struct{})
	go func() {
		r.shutdownWait.Wait()
		close(shutdone)
	}()
	select {
	case <-time.After(time.Second * 10):
	case <-shutdone:
	}
}

type WriteCloser struct {
	io.Writer
}

func (w *WriteCloser) Close() error { return nil }

func nopWriteCloser(w io.Writer) io.WriteCloser {
	return &WriteCloser{w}
}

func init() {
	globalhandle = make(map[string]Handle)
}
