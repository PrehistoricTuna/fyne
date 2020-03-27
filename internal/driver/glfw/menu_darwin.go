// +build !no_native_menus

package glfw

import (
	"unsafe"

	"fyne.io/fyne"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework AppKit

#include <AppKit/AppKit.h>

// Using void* as type for pointers is a workaround. See https://github.com/golang/go/issues/12065.
const void* darwinAppMenu();
const void* createDarwinMenu(const char* label);
void insertDarwinMenuItem(const void* menu, const char* label, int id, int index, bool separate);
void completeDarwinMenu(void* menu);
*/
import "C"

var callbacks []func()

//export menu_callback
func menu_callback(id int) {
	callbacks[id]()
}

func hasNativeMenu() bool {
	return true
}

func setupNativeMenu(main *fyne.MainMenu) {
	nextItemID := 0
	for _, menu := range main.Items {
		nextItemID = addNativeMenu(menu, nextItemID)
	}
}

func addNativeMenu(menu *fyne.Menu, nextItemID int) int {
	createMenu := false
	for _, item := range menu.Items {
		if !item.PlaceInNativeMenu {
			createMenu = true
			break
		}
	}

	var nsMenu unsafe.Pointer
	if createMenu {
		nsMenu = C.createDarwinMenu(C.CString(menu.Label))
	}

	for _, item := range menu.Items {
		if item.PlaceInNativeMenu {
			C.insertDarwinMenuItem(
				C.darwinAppMenu(),
				C.CString(item.Label),
				C.int(nextItemID),
				C.int(1),
				C.bool(item.Separate),
			)
		} else {
			C.insertDarwinMenuItem(
				nsMenu,
				C.CString(item.Label),
				C.int(nextItemID),
				C.int(-1),
				C.bool(item.Separate),
			)
		}
		callbacks = append(callbacks, item.Action)
		nextItemID++
	}

	if nsMenu != nil {
		C.completeDarwinMenu(nsMenu)
	}
	return nextItemID
}
