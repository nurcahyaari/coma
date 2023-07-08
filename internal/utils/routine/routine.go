package routine

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Routine struct {
	wg      *sync.WaitGroup
	mtx     *sync.Mutex
	job     map[string]Job
	results map[string]any
	params  map[string][]any
	err     map[string]error
	counter int
	maxPool int
}

func New() *Routine {
	return &Routine{
		wg:      &sync.WaitGroup{},
		mtx:     &sync.Mutex{},
		job:     make(map[string]Job),
		err:     make(map[string]error),
		params:  make(map[string][]any),
		results: make(map[string]any),
	}
}

type Job func(params ...any) (any, error)

// define maximum worker pool
func (r *Routine) SetMaximumPool(maxPool int) {
	r.maxPool = maxPool
}

// adding jobs
func (r *Routine) Add(name string, result any, job Job, params ...any) error {
	if _, exists := r.job[name]; exists {
		return errors.New("err: routine name is duplicated")
	}
	r.mtx.Lock()
	r.job[name] = job
	r.params[name] = append(r.params[name], params...)
	r.results[name] = result
	r.counter++
	r.mtx.Unlock()
	return nil
}

// start the job
func (r *Routine) Start() {
	r.wg.Add(r.counter)
	concurencyLimit := make(chan int, r.maxPool)
	for name, job := range r.job {
		if r.maxPool > 0 {
			concurencyLimit <- 1
		}
		func() {
			r.run(name, job)
			if r.maxPool > 0 {
				<-concurencyLimit
			}
		}()

	}

	r.wg.Wait()
}

func (r *Routine) run(name string, job Job) {
	defer r.wg.Done()
	res, err := job(r.params[name]...)

	r.mtx.Lock()

	sourceResult := reflect.ValueOf(res)
	if sourceResult.Kind() != reflect.Ptr {
		r.err[name] = errors.New("result must pass a pointer")
		return
	}

	destResult := reflect.ValueOf(r.results[name])
	if destResult.Kind() != reflect.Ptr {
		r.err[name] = errors.New("destination must pass a pointer")
		return
	}

	// Check if the types are assignable
	sourceType := sourceResult.Type()
	destType := destResult.Type()
	if !sourceType.AssignableTo(destType) {
		r.err[name] = fmt.Errorf("err: destination is unassignable from the result. Source type: %s, Destination type: %s",
			sourceType.String(),
			destType.String(),
		)
		return
	}

	// Set the value of the destination
	destResult.Elem().Set(sourceResult.Elem())

	if err != nil {
		r.err[name] = err
	}
	r.mtx.Unlock()
}

// check error
func (r *Routine) Errors() []error {
	errs := []error{}
	for name, err := range r.err {
		errs = append(errs, fmt.Errorf("%s: %s", name, err))
	}
	return errs
}

func (r *Routine) Error() error {
	errs := []string{}
	for name, err := range r.err {
		errs = append(errs, fmt.Sprintf("%s: %s", name, err))
	}
	return errors.New(strings.Join(errs, ";"))
}

func (r *Routine) IsError() bool {
	return len(r.err) > 0
}
