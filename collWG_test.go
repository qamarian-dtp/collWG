package collWG

import (
	"fmt"
	"gopkg.in/qamarian-lib/err.v0" // v0.4.0
	"gopkg.in/qamarian-lib/str.v3"
	"sync"
	"sync/atomic"
	"testing"
)

func TestCollWG (t *testing.T) {
	str.PrintEtr ("Test started.", "nte", "TestCollWG ()")

	var waitN uint16 = (1 << 16) - 1
	cWG, cWGP := New (waitN)

	goWG := &sync.WaitGroup {}
	goWG.Add (int (waitN) + 4)

	var okys int32 = 0
	var errs int32 = 0

	for i := 1; i <= int (waitN) + 4; i = i + 1 {
		go func (j int) {
			errX := cWGP.Done ()
			if errX != nil {
				if errX != ErrDone {
					errMssg := fmt.Sprintf ("Error at iteration %d. " +
						"[%s]", j, err.Fup (errX))
					str.PrintEtr (errMssg, "err", "TestCollWG ()")
					t.Fail ()
				} else {
					// str.PrintEtr ("Done!", "nte", "TestCollWG ()")
					atomic.AddInt32 (&errs, 1)
				}
			} else {
				atomic.AddInt32 (&okys, 1)
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
	cWG.Wait ()
	str.PrintEtr ("Test passed.", "std", "TestCollWG ()")
}
