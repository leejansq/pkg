package enginer

import (
	//"bytes"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/leejansq/pkg/log"
)

type Job struct {
	Eng      *Enginer
	Name     string
	Env      map[string]interface{}
	handler  Handle
	Args     []string
	Stdout   *Output
	Stderr   *Output
	Stdin    *Input
	closeIO  bool
	canclled chan struct{}
	once     sync.Once
	end      time.Time
}

func (job *Job) Run() (err error) {
	defer func() {
		// Wait for all background tasks to complete
		if job.closeIO {
			if err := job.Stdout.Close(); err != nil {
				log.Error(err)
			}
			if err := job.Stderr.Close(); err != nil {
				log.Error(err)
			}
			if err := job.Stdin.Close(); err != nil {
				log.Error(err)
			}
		}
	}()

	job.Eng.l.Lock()
	job.Eng.tasks.Add(1)
	job.Eng.l.Unlock()
	defer job.Eng.tasks.Done()

	if job.Eng.logging {
		log.Info("+job", job.CallString())
		defer func() {
			okerr := "OK"
			if err != nil {
				okerr = fmt.Sprintf("ERR: %s", err)
			}
			log.Info("-job", job.CallString(), okerr)
		}()
	}

	if job.handler == nil {
		return fmt.Errorf("%s: command not found", job.Name)
	}

	//var errorMessage = bytes.NewBuffer(nil)
	//job.Stderr.Add(errorMessage)

	err = job.handler(job)
	job.end = time.Now()

	return
}

func (job *Job) CallString() string {
	return fmt.Sprintf("%s(%s)", job.Name, strings.Join(job.Args, ", "))
}

func (job *Job) SetCloseIO(val bool) {
	job.closeIO = val
}

func (job *Job) Cancel() {
	job.once.Do(func() {
		close(job.canclled)
	})
}

func (job *Job) WaitCancelled() <-chan struct{} {
	return job.canclled
}

func (job *Job) SetEnv(name string, val interface{}) {
	job.Env[name] = val
}

func (job *Job) GetEnv(name string) (interface{}, bool) {
	val, ok := job.Env[name]
	return val, ok
}

func (job *Job) GetEnvString(name string) string {
	val, ok := job.Env[name]
	if ok {
		switch val.(type) {
		case string:
			return val.(string)
		}
	}
	return ""
}
func (job *Job) GetEnvInt(name string) int {
	val, ok := job.Env[name]
	if ok {
		switch val.(type) {
		case int:
			return val.(int)
		}
	}
	return -1
}
