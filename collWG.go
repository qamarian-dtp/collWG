package collWG

import (
	"fmt"
	"gopkg.in/qamarian-dtp/err.v0" // v0.4.0
	"sync"
)

func New (x uint16) (PvfCollWG, PbfCollWG) {
	w := &waitGroup {&sync.WaitGroup {}, x, &sync.Mutex {}}
	w.w.Add (int (x))
	return w, w
}

type PvfCollWG interface {
	Wait ()
}

type PbfCollWG interface {
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
	if w.x == 0 {
		return ErrDone
	}
	// --1-- ]

	// --1-- [
	w.w.Done ()
	w.x = w.x - 1
	// --1-- ]
	
	return nil
}

func (w *waitGroup) Wait () {
	w.w.Wait ()
}

var (
	ErrDone error = err.New ("Waiter is done waiting.", nil, nil)
)
