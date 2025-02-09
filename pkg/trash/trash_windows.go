//go:build windows

// Package trash
//
// Got inspiration from the following articles:
// - https://www.codeproject.com/Articles/2783/How-to-programmatically-use-the-Recycle-Bin
// - https://justen.codes/breaking-all-the-rules-using-go-to-call-windows-api-2cbfd8c79724
package trash

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"golang.org/x/sys/windows"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"
)

type shFileOpStruct struct {
	hwnd                  uintptr
	wFunc                 uintptr
	pFrom                 uintptr
	pTo                   uintptr
	fileOpFlags           uintptr
	fAnyOperationsAborted uintptr
	hNameMappings         uintptr
	lpszProgressTitle     uintptr
}

const (
	FO_DELETE          = 0x3
	FOF_ALLOWUNDO      = 0x40
	FOF_NOCONFIRMATION = 0x10
	FOF_NOERRORUI      = 0x400
	FOF_SILENT         = 0x4
)

var (
	shell32             = syscall.NewLazyDLL("shell32.dll")
	procSHFileOperation = shell32.NewProc("SHFileOperationW")
)

func Put(filenames ...string) error {
	pFromData := make([]uint16, 0, 256)
	for _, filename := range filenames {
		absPath, err := filepath.Abs(filename)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %v", err)
		}

		ptr, err := windows.UTF16FromString(absPath)
		if err != nil {
			return fmt.Errorf("failed to convert path %q: %v", absPath, err)
		}
		pFromData = append(pFromData, ptr...)
	}
	pFromData = append(pFromData, 0)

	title := []uint16{0, 0}

	param := &shFileOpStruct{
		wFunc:             FO_DELETE,
		pFrom:             uintptr(unsafe.Pointer(&pFromData[0])),
		fileOpFlags:       FOF_ALLOWUNDO | FOF_NOCONFIRMATION | FOF_SILENT,
		lpszProgressTitle: uintptr(unsafe.Pointer(&title[0])),
	}

	ret, _, _ := procSHFileOperation.Call(uintptr(unsafe.Pointer(param)))
	if ret != 0 {
		return fmt.Errorf("operation on %s failed with error code: %v", filenames, ret)
	}

	return nil
}

func Restore(fileName string) error {
	// Initialize COM
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	// Get the Shell.Application COM object
	unknown, err := oleutil.CreateObject("Shell.Application")
	if err != nil {
		return fmt.Errorf("failed to create Shell.Application object: %v", err)
	}
	defer unknown.Release()

	shellApp, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return fmt.Errorf("failed to get IDispatch interface: %v", err)
	}
	defer shellApp.Release()

	// Get the Recycle Bin folder
	recycleBinFolder := oleutil.MustCallMethod(shellApp, "NameSpace", 10).ToIDispatch()
	if recycleBinFolder == nil {
		return fmt.Errorf("failed to get Recycle Bin folder")
	}
	defer recycleBinFolder.Release()

	// Enumerate items in the Recycle Bin
	items := oleutil.MustCallMethod(recycleBinFolder, "Items").ToIDispatch()
	if items == nil {
		return fmt.Errorf("failed to enumerate items in Recycle Bin")
	}
	defer items.Release()

	count := oleutil.MustGetProperty(items, "Count").Val
	for i := 0; i < int(count); i++ {
		item := oleutil.MustCallMethod(items, "Item", i).ToIDispatch()
		if item == nil {
			continue
		}

		name := oleutil.MustGetProperty(item, "Name").ToString()
		if name == fileName {
			_, err = item.CallMethod("InvokeVerb", "undelete")
			if err != nil {
				return err
			}
			// InvokeVerb is asynchronous so we need to wait for the call to finish.
			for range 30 {
				originalPath := getOriginalPath(recycleBinFolder, item)
				if originalPath != "" {
					break
				}
				time.Sleep(10 * time.Millisecond)
			}
			item.Release()
			return nil
		}
		item.Release()
	}

	return fmt.Errorf("file '%s' not found in Recycle Bin", fileName)
}

func getOriginalPath(folder *ole.IDispatch, item *ole.IDispatch) string {
	// The column index for "Original Location" varies by system and language.
	// It's usually the second column, so index 1. You may need to adjust this.
	const originalPathColumnIndex = 1
	details := oleutil.MustCallMethod(folder, "GetDetailsOf", item, originalPathColumnIndex).ToString()
	return details
}
