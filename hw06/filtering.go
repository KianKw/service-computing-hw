package rxgo

import (
    "context"
    "errors"
	"reflect"
	"sync"
	"time"
	"fmt"
)

var ErrInput = errors.New("Input Error!");

var ErrIndex = errors.New("The index is illegal!")

type filteringOperator struct {
	opFunc func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool)
}

// initialize a new FilterObservable
func (parent *Observable) newFilterObservable(name string) (o *Observable) {
    o = newObservable()
    o.Name = name
    o.root = parent.root
    o.pred = parent
    parent.next = o
    o.buf_len = BufferLen
    o.only_first = false
    o.only_last = false
    o.debounce_timespan = 0
    o.only_distinct = false
    o.skip = 0
    o.take = 0
    return
}

// Debounce only emit an item from an Observable if a particular timespan has passed without it emitting another item
func (parent *Observable) Debounce(t time.Duration) (o *Observable) {
	o = parent.newFilterObservable("Debounce")
	o.debounce_timespan = t
    o.operator = filteringOperator {
		opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
			var params = []reflect.Value{x}
			x = params[0]
			if !end {
				end = o.sendToFlow(ctx, x.Interface(), out)
			}
			return
		},
	}
	return o
}

// Distinct suppress duplicate items emitted by an Observable
func (parent *Observable) Distinct() (o *Observable) {
	o = parent.newFilterObservable("Distinct")
	o.only_distinct = true
    o.operator = filteringOperator {
		opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
			var params = []reflect.Value{x}
			x = params[0]
			if !end {
				end = o.sendToFlow(ctx, x.Interface(), out)
			}
			return
		},
	}
	return o
}

// ElementAt emit only item n emitted by an Observable
func (parent *Observable) ElementAt(id int) (o *Observable) {
	o = parent.newFilterObservable("ElementAt.n")
	o.element_at = id
    o.operator = filteringOperator {
		opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
			var params = []reflect.Value{x}
			x = params[0]
			if !end {
				end = o.sendToFlow(ctx, x.Interface(), out)
			}
			return
		},
	}
	return
}

// First emit only the first item (or the first item that meets some condition) emitted by an Observable
func (parent *Observable) First() (o *Observable) {
	o = parent.newFilterObservable("First")
	o.only_first = true
    o.operator = filteringOperator {
		opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
			var params = []reflect.Value{x}
			x = params[0]
			if !end {
				end = o.sendToFlow(ctx, x.Interface(), out)
			}
			return
		},
	}
	return o
}

// Last emit only the last item (or the last item that meets some condition) emitted by an Observable
func (parent *Observable) Last() (o *Observable) {
	o = parent.newFilterObservable("Last")
	o.only_last  = true
    o.operator = filteringOperator {
		opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
			var params = []reflect.Value{x}
			x = params[0]
			if !end {
				end = o.sendToFlow(ctx, x.Interface(), out)
			}
			return
		},
	}
	return o
}

// Sample emit the most recent item emitted by an Observable within periodic time intervals.
func (parent *Observable) Sample(timespan time.Duration) (o *Observable) {
	o = parent.newFilterObservable("Sample")
	o.sample_interval = timespan
    o.element_at = 0
    o.operator = filteringOperator {
		opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
			var params = []reflect.Value{x}
			x = params[0]
			if !end {
				end = o.sendToFlow(ctx, x.Interface(), out)
			}
			return
		},
	}
	return o
}

// Skip suppress the first n items emitted by an Observable
func (parent *Observable) Skip(num int) (o *Observable) {
	o = parent.newFilterObservable("Skip.n")
	o.skip = num
    o.operator = filteringOperator {
		opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
			var params = []reflect.Value{x}
			x = params[0]
			if !end {
				end = o.sendToFlow(ctx, x.Interface(), out)
			}
			return
		},
	}
	return o
}

// SkipLast suppress the last n items emitted by an Observable
func (parent *Observable) SkipLast(num int) (o *Observable) {
	o = parent.newFilterObservable("SkipLast.n")
	o.skip = -num
    o.operator = filteringOperator {
		opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
			var params = []reflect.Value{x}
			x = params[0]
			if !end {
				end = o.sendToFlow(ctx, x.Interface(), out)
			}
			return
		},
	}
	return o
}

// Take emit only the first n items emitted by an Observable
func (parent *Observable) Take(num int) (o *Observable) {
	o = parent.newFilterObservable("Take")
	o.take = num
    o.is_taking = true
    o.operator = filteringOperator {
		opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
			var params = []reflect.Value{x}
			x = params[0]
			if !end {
				end = o.sendToFlow(ctx, x.Interface(), out)
			}
			return
		},
	}
	return o
}

// TakeLast emit only the last n items emitted by an Observable
func (parent *Observable) TakeLast(num int) (o *Observable) {
	o = parent.newFilterObservable("TakeLast")
	o.take = -num
    o.is_taking = true
    o.operator = filteringOperator {
		opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
			var params = []reflect.Value{x}
			x = params[0]
			if !end {
				end = o.sendToFlow(ctx, x.Interface(), out)
			}
			return
		},
	}
	return o
}

func (tsop filteringOperator) op(ctx context.Context, o *Observable) {
	in := o.pred.outflow
    out := o.outflow
    interval := o.debounce_timespan
    var wg sync.WaitGroup
    var out_buf []interface{}
	go func() {
		is_appear := make(map[interface{}]bool)
		end := false
		start := time.Now()
		sample_start := time.Now()
		for x := range in {
			rt := time.Since(start)
			st := time.Since(sample_start)
			start = time.Now()
			if end {
				continue
			}
			if o.sample_interval > 0 && st < o.sample_interval {
				continue
			}
			if interval > time.Duration(0) && rt < interval {
				continue
			}
			xv := reflect.ValueOf(x)
			if e, ok := x.(error); ok && !o.flip_accept_error {
				o.sendToFlow(ctx, e, out)
				continue
			}
			o.mu.Lock()
			out_buf = append(out_buf, x)
			o.mu.Unlock()
			if o.element_at > 0 {
				continue
			}
			if o.take != 0 || o.skip != 0 {
				continue
			}
			if o.only_last {
				continue
			}
			if o.only_distinct && is_appear[xv.Interface()] {
				continue
			}
			o.mu.Lock()
			is_appear[xv.Interface()] = true
			o.mu.Unlock()
			if interval > time.Duration(0) {
				if len(out_buf) > 2 {
					xv = reflect.ValueOf(out_buf[len(out_buf) - 2])
				}
			}
			switch threading := o.threading; threading {
			case ThreadingDefault:
				if o.sample_interval > 0 {
					sample_start = sample_start.Add(o.sample_interval)
				}
				if tsop.opFunc(ctx, o, xv, out) {
					end = true
				}
			case ThreadingIO:
				fallthrough
			case ThreadingComputing:
				wg.Add(1)
				if o.sample_interval > 0 {
					sample_start.Add(o.sample_interval)
				}

				go func() {
					defer wg.Done()
					if tsop.opFunc(ctx, o, xv, out) {
						end = true
					}
				}()
			default:
			}
			if o.only_first {
				break
			}
		}
		if o.only_last && len(out_buf) > 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				xv := reflect.ValueOf(out_buf[len(out_buf)-1])
				tsop.opFunc(ctx, o, xv, out)
			}()
		}
		if o.take != 0 || o.skip != 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var div int
				if o.is_taking {
					div = o.take
				} else {
					div = o.skip
				}
				new_in,err := judgeMode(o.is_taking, div, out_buf)

				if err != nil {
					o.sendToFlow(ctx, err, out)
				} else {
					xv := new_in
					for _, val := range xv {
						tsop.opFunc(ctx, o, reflect.ValueOf(val), out)
					}
				}
			}()
		}
		if o.element_at != 0 {
			if o.element_at < 0 || o.element_at > len(out_buf) {
				o.sendToFlow(ctx, ErrIndex, out)
			} else {
				xv := reflect.ValueOf(out_buf[o.element_at-1])
				tsop.opFunc(ctx, o, xv, out)
			}
		}
		wg.Wait()
		if (o.only_last || o.only_first) && len(out_buf) == 0 && !o.flip_accept_error  {
			o.sendToFlow(ctx, ErrInput, out)
		}
		o.closeFlow(out)
	}()
}

// If isTake is true, then switch to Take mode, otherwise switch to Skip mode
func judgeMode(isTake bool, division int, in []interface{}) ([]interface{}, error) {
	fmt.Println(in)
	if (isTake && division > 0) || (!isTake && division < 0) {
		if !isTake {
			division = len(in) + division
		}
		if division >= len(in) || division <= 0{
			return nil, ErrIndex
		}
		return in[:division], nil
	}

	if (isTake && division < 0) || (!isTake && division > 0) {
		if isTake {
			division = len(in) + division
		}
		if division >= len(in) || division <= 0{
			return nil, ErrIndex
		}
		return in[division:], nil
	}
	return nil, ErrIndex
}
