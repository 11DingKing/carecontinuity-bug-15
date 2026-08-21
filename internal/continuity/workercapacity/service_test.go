package workercapacity_test

import (
	"carecontinuity/internal/continuity/workercapacity"
	"sync"
	"sync/atomic"
	"testing"
)

func TestCommunityWorkerCapacityPublicBehavior(t *testing.T) {
	c := workercapacity.NewCoordinator()
	release := make(chan struct{})
	var arrived atomic.Int32
	c.SetBarrier(func() {
		if arrived.Add(1) == 2 {
			close(release)
		}
		<-release
	})
	var wg sync.WaitGroup
	out := make(chan bool, 2)
	for _, task := range []string{"testing", "followup"} {
		wg.Add(1)
		go func(v string) { defer wg.Done(); out <- c.Assign("team-river", v) }(task)
	}
	wg.Wait()
	close(out)
	accepted := 0
	for ok := range out {
		if ok {
			accepted++
		}
	}
	if accepted != 1 {
		t.Fatalf("expected one team assignment, got %d", accepted)
	}
	if owner, ok := c.Owner("team-river"); !ok || (owner != "testing" && owner != "followup") {
		t.Fatalf("invalid final owner %q %v", owner, ok)
	}
}
