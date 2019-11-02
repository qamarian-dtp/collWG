package cwg

import (
	"fmt"
	"gopkg.in/qamarian-dtp/err.v0" // v0.4.0
	"sync"
)

// New () creates a new collected wait group. Argument should be the number of goroutines
// to wait for.
func New (x uint16) (PvfCWG, PbfCWG) {
	w := &waitGroup {&sync.WaitGroup {}, x, &sync.Mutex {}}
	w.w.Add (int (x))
	return w, w
}

// PvfCollWG is the private face of the collected wait group.
// This data is not necessarily thread-safe.
type PvfCWG interface {
	Wait ()
}

// PvfCollWG is the public face of the collected wait group.
// This data is thread-safe.
type PbfCWG interface {
	// Possible errors include: ErrDone.
	Done () (error)
}

type waitGroup struct {
	w *sync.WaitGroup // Underlying wait group.
	x uint16          // No of rountines to wait for.
	l *sync.Mutex     // Wait group lock.
}

func (w *waitGroup) Done () (e error) {
	defer func () {
		panicReason := recover ()
		if panicReason != nil {
			panicErr := err.New (fmt.Sprintf ("%v", panicReason), nil, nil)
			e = err.New ("A panic occured.", nil, nil, panicErr)
		}
	} ()

	// --1-- [	
	if w.x == 0 {
		return ErrDone
	}
	// --1-- ]

	// --1-- [
	w.l.Lock ()
	defer w.l.Unlock ()
	
	// Checking again, just in case value had become 0, while trying to lock data.
	if w.x == 0 {
		return ErrDone
	}
	// --1-- ]

	// --1-- [
	w.w.Done ()
	w.x --
	// --1-- ]
	
	return nil
}

func (w *waitGroup) Wait () {
	w.w.Wait ()
}

var (
	ErrDone error = err.New ("Waiter is done waiting.", nil, nil)
)
