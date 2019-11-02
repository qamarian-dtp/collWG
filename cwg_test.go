package cwg

import (
	"fmt"
	"gopkg.in/qamarian-lib/err.v0" // v0.4.0
	"gopkg.in/qamarian-lib/str.v3"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestCollWG () extensively tests the "collected wait group" data type.
func TestCollWG (t *testing.T) {
	str.PrintEtr ("Test started.", "std", "TestCollWG ()")

	// Testing method Done (). --1-- [
	// | --
	var waitN uint16 = (1 << 16) - 1
	var additional uint8 = 4
	
	_, pbFace := New (waitN)

	goWG := &sync.WaitGroup {}
	goWG.Add (int (waitN) + int (additional))

	var okys int32 = 0
	var errs int32 = 0
	// -- |

	for i := 1; i <= int (waitN) + int (additional); i = i + 1 {
		go func (j int) {
			errX := pbFace.Done ()
			if errX == nil {
				atomic.AddInt32 (&okys, 1)
			} else {
				if errX == ErrDone {
					atomic.AddInt32 (&errs, 1)
				} else {
					errMssg := fmt.Sprintf ("Error at iteration %d. " +
						"[%s]", j, err.Fup (errX))
					str.PrintEtr (errMssg, "err", "TestCollWG ()")
					t.Fail ()
				}
			}
			goWG.Done ()
		} (i)
	}

	goWG.Wait ()

	if t.Failed () {
		str.PrintEtr ("Test failed. Ref: 0", "err", "TestCollWG ()")
		t.FailNow ()
	}
	if okys != int32 (waitN) || errs != int32 (additional) {
		str.PrintEtr ("Test failed. Ref: 1", "err", "TestCollWG ()")
		t.FailNow ()
	}
	// --1-- ]

	// Testing relationship between method Wait () and Done (). --1-- [
	// | --
	var waitN2 uint16 = (1 << 16) - 1

	pvFace2, pbFace2 := New (waitN2)
	
	goWG2 := &sync.WaitGroup {}
	goWG2.Add (int(waitN2) - 1)
	
	done := false
	doneLock := &sync.Mutex {}
	// -- |
	
	for i := 1; i < int (waitN2); i ++ {
		go func () {
			pbFace2.Done ()
			goWG2.Done ()
		} ()
	}
	go func () {
		goWG2.Wait ()
		time.Sleep (time.Second * 4)
		doneLock.Lock ()
		defer doneLock.Unlock ()
		done = true
		pbFace2.Done ()
	} ()
		
	pvFace2.Wait ()

	doneLock.Lock ()
	defer doneLock.Unlock ()
	
	if done == false {
		str.PrintEtr ("Test failed. Ref: 2", "err", "TestCollWG ()")
		t.FailNow ()
	}
	// --1-- ]	

	str.PrintEtr ("Test passed.", "std", "TestCollWG ()")
}
