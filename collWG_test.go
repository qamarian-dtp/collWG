package collWG

import (
	"fmt"
	"gopkg.in/qamarian-lib/err.v0" // v0.4.0
	"gopkg.in/qamarian-lib/str.v3"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCollWG (t *testing.T) {
	str.PrintEtr ("Test started.", "std", "TestCollWG ()")

	// Testing method Done (). --1-- [
	// | --
	var waitN uint16 = (1 << 16) - 1

	_, clWGP := New (waitN)

	goWG := &sync.WaitGroup {}
	goWG.Add (int (waitN) + 4)

	var okys int32 = 0
	var errs int32 = 0
	// -- |

	for i := 1; i <= int (waitN) + 4; i = i + 1 {
		go func (j int) {
			errX := clWGP.Done ()
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
	if okys != int32 (waitN) || errs != 4 {
		str.PrintEtr ("Test failed. Ref: 1", "err", "TestCollWG ()")
		t.FailNow ()
	}
	// --1-- ]

	// Testing relationship between method Wait () and Done (). --1-- [
	// | --
	var waitN2 uint16 = (1 << 16) - 1

	clWG2, clWGP2 := New (waitN2)
	
	goWG2 := &sync.WaitGroup {}
	goWG2.Add (int(waitN2) - 1)
	
	done := false
	// -- |
	
	for i := 1; i < int (waitN2); i ++ {
		go func () {
			clWGP2.Done ()
			goWG2.Done ()
		} ()
	}
	go func () {
		goWG2.Wait ()
		time.Sleep (time.Second * 4)
		done = true
		clWGP2.Done ()
	} ()
		
	clWG2.Wait ()
	
	if done == false {
		str.PrintEtr ("Test failed. Ref: 2", "err", "TestCollWG ()")
		t.FailNow ()
	}

	time.Sleep (time.Second * 1)
	// --1-- ]	

	str.PrintEtr ("Test passed.", "std", "TestCollWG ()")
}
