# RxGo: Reactive Extensions for the Go

> This repository is fork form [pmlpml/rxgo](https://github.com/pmlpml/rxgo)

RxGo is a Go implementation of [Reactive Extensions](http://reactivex.io/documentation/observable.html): a library for composing asynchronous and event-based programs by using observable sequences.

It extends the observer pattern to support sequences of data/events and adds operators that allow you to compose sequences together declaratively while abstracting away concerns about things like low-level threading, synchronization, thread-safety and concurrent data structures.

#### V0.2

* function program friendly
* Connectable observables, support Hot or Cold item emiting
* context support for concurrency
* support error type data in observable sequences
* internal debug support
* unsubscribe supporting
* non-opinionated about source of concurrency 
* async or synchronous execution

## Getting started

### Hello World

The **Hello World** program:

```go
package main

import (
	"fmt"
	RxGo "github.com/pmlpml/rxgo"
)

func main() {
	RxGo.Just("Hello", "World", "!").Subscribe(func(x string) {
		fmt.Println(x)
	})
}
```

output:

```
Hello
World
!
```

### Chained operations on stream

The dataflows in RxGo consist of a source, zero or more intermediate steps followed by a data consumer or combinator step (where the step is responsible to consume the dataflow by some means):

```
source.operator1().operator2().operator3().subscribe(observer);
```

The first obserable emits items, so called _source_. And a chained _operations_ apply on each item asynchronously, until subscribe to the _observer_.

The sample program generate a **fibonacci sequence**, and then double each item.

```go
package main

import (
	"fmt"
	RxGo "github.com/pmlpml/rxgo"
)

func fibonacci(max int) func() (int, bool) {
	a, b := 0, 1
	return func() (r int, end bool) {
		r = a
		a, b = b, a+b
		if r > max {
			end = true
		}
		return
	}
}

func main() {
	RxGo.Start(fibonacci(10)).Map(func(x int) int {
		return 2*x
	}).Subscribe(func(x int) {
		fmt.Print(x)
	})
}
```

The result is `022461016`.  the source `Start(fibonacci(10))` generates dataflow `0112358` to a operation `Map`, and `Subscribe` to print

### Connectable observables

A Connectable Observable resembles an ordinary Observable, except that it does not begin emitting items when it is subscribed to, but only when its connect() method is called. 
There are two phases for data pipeline.

The first phase is called _definition_. we define source function or operators for a pipeline that is the same as assign work roles for each node on it

The next phase is called _runtime_. when any observable call `Subscribe(...)` , the pipeline include this observable will connect its all node or 
appoint one or more workers stand aside by each node. if predecess worker emit a data, the worker will play the role defined before, 
and send result to the next worker

The simple program shows the function of a **Pipeline Restart**

```go
package main

import (
	"fmt"
	RxGo "github.com/pmlpml/rxgo"
)

func main() {
    //define pipeline
	source := RxGo.Just("Hello", "World", "!")
	next := source.Map(func(s string) string {
		return s + " "
	})
    //run pipeline
	next.Subscribe(func(x string) {
		fmt.Print(x)
	})
	fmt.Println()
	source.Subscribe(func(x string) {
		fmt.Print(x)
	})
}
```

the program will print `Hello World ! ` twice!


