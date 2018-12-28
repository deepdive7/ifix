package ifix

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func no() bool  { return false }
func yes() bool { return true }

type A struct {
	X int
}

func (a *A) Show() {
	fmt.Println("X = ", a.X)
}

func (a *A) Play() {
	fmt.Println(-1)
}

func Show(a *A) {
	a.Show()
}

func Hey() {
	fmt.Println(-1)
}

func TestMultiplePatch(t *testing.T) {
	a := &A{}
	for i := 1; i < 8; i++ {
		x := i
		PatchInstanceMethod(reflect.TypeOf(a), "Play", func(a *A) {
			fmt.Fprintln(os.Stdout, "X = ", x)
		})
		Patch(Hey, func() {
			fmt.Println(x*x)
		})
		a.Play()
		Hey()
	}
}

func TestA(t *testing.T) {
	a := &A{3}
	fmt.Print("old a.Show->") //old a.Show->X =  3
	Show(a)
	Patch((*A).Show, func(a *A) {
		fmt.Fprintln(os.Stdout, "X^2 = ", a.X*a.X)
	})
	b := &A{2}
	fmt.Print("new b.Show->") //new b.Show->X^2 =  4
	Show(b)
	fmt.Print("new a.Show->") //new a.Show->X^2 =  9
	Show(a)
}

func TestTimePatch(t *testing.T) {
	before := time.Now()
	Patch(time.Now, func() time.Time {
		return time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	})
	during := time.Now()
	assert.True(t, Unpatch(time.Now))
	after := time.Now()

	assert.Equal(t, time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), during)
	assert.NotEqual(t, before, during)
	assert.NotEqual(t, during, after)
}

func TestGC(t *testing.T) {
	value := true
	Patch(no, func() bool {
		return value
	})
	defer UnpatchAll()
	runtime.GC()
	assert.True(t, no())
}

func TestSimple(t *testing.T) {
	assert.False(t, no())
	Patch(no, yes)
	assert.True(t, no())
	assert.True(t, Unpatch(no))
	assert.False(t, no())
	assert.False(t, Unpatch(no))
}

func TestGuard(t *testing.T) {
	var guard *PatchGuard
	guard = Patch(no, func() bool {
		guard.Unpatch()
		defer guard.Restore()
		return !no()
	})
	for i := 0; i < 100; i++ {
		assert.True(t, no())
	}
	Unpatch(no)
}

func TestUnpatchAll(t *testing.T) {
	assert.False(t, no())
	Patch(no, yes)
	assert.True(t, no())
	UnpatchAll()
	assert.False(t, no())
}

type s struct{}

func (s *s) yes() bool { return true }

func TestWithInstanceMethod(t *testing.T) {
	i := &s{}

	assert.False(t, no())
	Patch(no, i.yes)
	assert.True(t, no())
	Unpatch(no)
	assert.False(t, no())
}

type f struct{}

func (f *f) No() bool { return false }

func TestOnInstanceMethod(t *testing.T) {
	i := &f{}
	assert.False(t, i.No())
	PatchInstanceMethod(reflect.TypeOf(i), "No", func(_ *f) bool { return true })
	assert.True(t, i.No())
	assert.True(t, UnpatchInstanceMethod(reflect.TypeOf(i), "No"))
	assert.False(t, i.No())
}

func TestNotFunction(t *testing.T) {
	assert.Panics(t, func() {
		Patch(no, 1)
	})
	assert.Panics(t, func() {
		Patch(1, yes)
	})
}

func TestNotCompatible(t *testing.T) {
	assert.Panics(t, func() {
		Patch(no, func() {})
	})
}
