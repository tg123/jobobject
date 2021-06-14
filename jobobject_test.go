// +build windows

package jobobject

import (
	"os/exec"
	"testing"
	"time"
)

func ExampleJobObject() {

	e := exec.Command("notepad.exe")
	e.Start()

	job, err := Create()
	if err != nil {
		panic(err)
	}

	err = job.AddProcess(e.Process)
	if err != nil {
		panic(err)
	}

	// You can either close it or not.
	// when the owner process of the job exits, all jobs will be closed
	defer job.Close()
}

func TestJobObject(t *testing.T) {
	e := exec.Command("notepad.exe")
	e.Start()

	job, err := Create()
	if err != nil {
		t.Fatal(err)
	}

	err = job.AddProcess(e.Process)
	if err != nil {
		t.Fatal(err)
	}

	closeCalled := false
	exited := false

	go func() {
		_, err := e.Process.Wait()
		if err != nil {
			t.Fatal(err)
		}

		if !closeCalled {
			t.Fatal("exit before close called")
		}

		exited = true
	}()

	time.Sleep(2 * time.Second) // let process live for 2s

	closeCalled = true
	// process may exit during this period :(
	if err := job.Close(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	if !exited {
		t.Fatal("process not killed")
	}
}
