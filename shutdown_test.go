package shutdown_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ztrue/shutdown"
)

func TestHooks(t *testing.T) {
	defer shutdown.Reset()

	data := map[string]int{}
	hook1 := func() {
		data["1:foo"]++
	}
	hook2 := func() {
		data["2:bar"]++
	}
	hook3 := func() {
		data["3:baz"]++
	}
	hook4 := func() {
		data["4:qux"]++
	}
	hook5 := func(sig os.Signal) {
		data["5:"+sig.String()]++
	}
	hook6 := func(sig os.Signal) {
		data["6:"+sig.String()]++
	}
	hook7 := func(sig os.Signal) {
		data["7:"+sig.String()]++
	}
	hook8 := func(sig os.Signal) {
		data["8:"+sig.String()]++
	}

	// Add hooks.
	shutdown.Add(hook1)
	key2 := shutdown.Add(hook2)
	key3 := "hook3"
	shutdown.AddWithKey(key3, hook3)
	key4 := "hook4"
	shutdown.AddWithKey(key4, hook4)
	shutdown.AddWithParam(hook5)
	key6 := shutdown.AddWithParam(hook6)
	key7 := "hook7"
	shutdown.AddWithKeyWithParam(key7, hook7)
	key8 := "hook8"
	shutdown.AddWithKeyWithParam(key8, hook8)

	hooks := shutdown.Hooks()
	if len(hooks) != 8 {
		t.Errorf(
			"len(shutdown.Hooks()) = %#v; want %#v",
			len(hooks), 8,
		)
	}

	// Assert data is not added yet.
	if len(data) != 0 {
		t.Errorf(
			"len(data) = %#v; want %#v",
			len(data), 0,
		)
	}
	// Run hooks.
	for _, hook := range hooks {
		hook(os.Interrupt)
	}

	dataKeys := []string{
		"1:foo",
		"2:bar",
		"3:baz",
		"4:qux",
		"5:interrupt",
		"6:interrupt",
		"7:interrupt",
		"8:interrupt",
	}
	for _, dataKey := range dataKeys {
		if data[dataKey] != 1 {
			t.Errorf(
				"data[%#v] = %#v; want %#v",
				dataKey, data[dataKey], 1,
			)
		}
	}
	// Assert no extra data keys.
	if len(data) != 8 {
		t.Errorf(
			"len(data) = %#v; want %#v",
			len(data), 8,
		)
	}

	// Reset data.
	data = map[string]int{}
	// Remove some hooks.
	shutdown.Remove(key2)
	shutdown.Remove(key4)
	shutdown.Remove(key6)
	shutdown.Remove(key8)

	hooks = shutdown.Hooks()
	if len(hooks) != 4 {
		t.Errorf(
			"len(shutdown.Hooks()) = %#v; want %#v",
			len(hooks), 4,
		)
	}

	// Assert data is not added yet.
	if len(data) != 0 {
		t.Errorf(
			"len(data) = %#v; want %#v",
			len(data), 0,
		)
	}
	// Run hooks.
	for _, hook := range hooks {
		hook(os.Interrupt)
	}

	dataKeys = []string{
		"1:foo",
		"3:baz",
		"5:interrupt",
		"7:interrupt",
	}
	for _, dataKey := range dataKeys {
		if data[dataKey] != 1 {
			t.Errorf(
				"data[%#v] = %#v; want %#v",
				dataKey, data[dataKey], 1,
			)
		}
	}
	// Assert no extra data keys.
	if len(data) != 4 {
		t.Errorf(
			"len(data) = %#v; want %#v",
			len(data), 4,
		)
	}
}

func TestListenSameProcess(t *testing.T) {
	defer shutdown.Reset()

	data := map[string]int{}
	var mutex sync.Mutex

	// Add 3 hooks.
	shutdown.Add(func() {
		mutex.Lock()
		defer mutex.Unlock()
		data["foo"]++
	})
	key := shutdown.Add(func() {
		mutex.Lock()
		defer mutex.Unlock()
		data["bar"]++
	})
	shutdown.Add(func() {
		mutex.Lock()
		defer mutex.Unlock()
		data["baz"]++
	})
	// Remove one of them.
	shutdown.Remove(key)

	go func() {
		// TODO Is there a better solution to make sure listening is started?
		time.Sleep(10 * time.Millisecond)
		p, err := os.FindProcess(os.Getpid())
		if err != nil {
			panic(err.Error())
		}
		err = p.Signal(os.Interrupt)
		if err != nil {
			panic(err.Error())
		}
	}()

	shutdown.Listen()

	if len(data) != 2 {
		t.Errorf(
			"len(data) = %#v; want %#v",
			len(data), 2,
		)
	}

	dataKeys := []string{"foo", "baz"}
	for _, dataKey := range dataKeys {
		if data[dataKey] != 1 {
			t.Errorf(
				"data[%#v] = %#v; want %#v",
				dataKey, data[dataKey], 1,
			)
		}
	}
}

func TestListenSeparateProcess(t *testing.T) {
	// TODO "Listen" coverage does not count as
	// far as code executed in a separate process,
	// no matter if it's actually tested and covered,
	// how to make it count?
	if os.Getenv("LISTEN") == "1" {
		// Add 3 hooks.
		shutdown.Add(func() {
			fmt.Println("foo")
		})
		key := shutdown.Add(func() {
			fmt.Println("bar")
		})
		shutdown.Add(func() {
			fmt.Println("baz")
		})
		// Remove one of them.
		shutdown.Remove(key)
		shutdown.Listen()
		return
	}

	var buf bytes.Buffer
	cmd := exec.Command(os.Args[0], "-test.run=TestListenSeparateProcess")
	cmd.Env = append(os.Environ(), "LISTEN=1")
	cmd.Stdout = &buf
	err := cmd.Start()
	if err != nil {
		panic(err.Error())
	}

	// TODO Better solution to wait for programm launch?
	time.Sleep(10 * time.Millisecond)

	err = cmd.Process.Signal(os.Interrupt)
	if err != nil {
		panic(err.Error())
	}
	err = cmd.Wait()
	if err != nil {
		panic(err.Error())
	}
	lines := strings.Split(buf.String(), "\n")
	if len(lines) < 3 {
		t.Errorf(
			"len(lines) = %#v; want >= %#v",
			len(lines), 3,
		)
	}

	fooFirst := lines[0] == "foo" && lines[1] == "baz"
	bazFirst := lines[0] == "baz" && lines[1] == "foo"
	if !fooFirst && !bazFirst {
		t.Errorf(
			"lines[0] = %#v; lines[1] = %#v; want %#v and %#v",
			lines[0], lines[1], "foo", "bar",
		)
	}

	if lines[2] != "PASS" {
		t.Errorf(
			"line[2] = %#v; want %#v",
			lines[2], "PASS",
		)
	}
}

func TestNew(t *testing.T) {
	data := map[string]int{}

	hook1 := func() {
		data["foo"]++
	}
	hook2 := func() {
		data["bar"]++
	}

	s1 := shutdown.New()
	s2 := shutdown.New()

	// Add hooks with the same key.
	s1.AddWithKey("hook", hook1)
	s2.AddWithKey("hook", hook2)

	if len(s1.Hooks()) != 1 {
		t.Errorf(
			"len(s1.Hooks()) = %#v; want %#v",
			len(s1.Hooks()), 1,
		)
	}
	if len(s2.Hooks()) != 1 {
		t.Errorf(
			"len(s2.Hooks()) = %#v; want %#v",
			len(s2.Hooks()), 1,
		)
	}

	// Run both hooks.
	s1.Hooks()["hook"](os.Interrupt)
	s2.Hooks()["hook"](os.Interrupt)

	// Make sure both hooks executed once.
	if len(data) != 2 {
		t.Errorf(
			"len(data) = %#v; want %#v",
			len(data), 2,
		)
	}
	dataKeys := []string{"foo", "bar"}
	for _, dataKey := range dataKeys {
		if data[dataKey] != 1 {
			t.Errorf(
				"data[%#v] = %#v; want %#v",
				dataKey, data[dataKey], 1,
			)
		}
	}
}
