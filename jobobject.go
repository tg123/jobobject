// +build windows

// Windows JobObject utils for kill all child processes when parent process exits
package jobobject

import (
	"os"
	"reflect"
	"unsafe"

	"golang.org/x/sys/windows"
)

const sizeofJobobjectExtendedLimitInformation = 144 // binary.size cannot handle uintptr

type JobObject struct {
	handle windows.Handle
}

func Create() (*JobObject, error) {
	handle, err := windows.CreateJobObject(nil, nil)

	if err != nil {
		return nil, err
	}

	extendedInfo := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}

	_, err = windows.SetInformationJobObject(handle, windows.JobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&extendedInfo)), sizeofJobobjectExtendedLimitInformation)

	if err != nil {
		return nil, err
	}

	obj := &JobObject{handle: handle}

	return obj, nil
}

func (j *JobObject) Close() error {
	return windows.CloseHandle(j.handle)
}

func (j *JobObject) AddProcess(p *os.Process) error {
	fv := reflect.ValueOf(p).Elem().FieldByName("handle")
	return windows.AssignProcessToJobObject(j.handle, windows.Handle(fv.Uint()))
}
