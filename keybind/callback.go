/*
    Provides a callback interface, very similar to that found in
    xevent/callback.go --- but only for key bindings.
*/
package keybind

import "code.google.com/p/jamslam-x-go-binding/xgb"
import "github.com/BurntSushi/xgbutil"
import "github.com/BurntSushi/xgbutil/xevent"

type KeyPressFun xevent.KeyPressFun

func connect(xu *xgbutil.XUtil, callback xgbutil.KeyBindCallback,
             evtype int, win xgb.Id, keyStr string) {
    // Get the mods/key first
    mods, keycode := ParseString(xu, keyStr)

    // Only do the grab if we haven't yet on this window.
    if xu.KeyBindGrabs(evtype, win, mods, keycode) == 0 {
        Grab(xu, win, mods, keycode)
    }

    // Always attach the callback.
    xu.AttachKeyBindCallback(evtype, win, mods, keycode, callback)

    // If we've never grabbed anything on this window before, we need to
    // make sure we can respond to it in the main event loop.
    var allCb xgbutil.Callback
    if evtype == xevent.KeyPress {
        allCb = xevent.KeyPressFun(RunKeyPressCallbacks)
    } else {
        allCb = xevent.KeyReleaseFun(RunKeyReleaseCallbacks)
    }
    if !xu.Connected(evtype, win, allCb) {
        allCb.Connect(xu, win)
    }
}

func (callback KeyPressFun) Connect(xu *xgbutil.XUtil, win xgb.Id,
                                    keyStr string) {
    connect(xu, callback, xevent.KeyPress, win, keyStr)
}

func (callback KeyPressFun) Run(xu *xgbutil.XUtil, event interface{}) {
    callback(xu, event.(xevent.KeyPressEvent))
}

type KeyReleaseFun xevent.KeyReleaseFun

func (callback KeyReleaseFun) Connect(xu *xgbutil.XUtil, win xgb.Id,
                                      keyStr string) {
    connect(xu, callback, xevent.KeyRelease, win, keyStr)
}

func (callback KeyReleaseFun) Run(xu *xgbutil.XUtil, event interface{}) {
    callback(xu, event.(xevent.KeyReleaseEvent))
}

// RunKeyPressCallbacks infers the window, keycode and modifiers from a
// KeyPressEvent and runs the corresponding callbacks.
func RunKeyPressCallbacks(xu *xgbutil.XUtil, ev xevent.KeyPressEvent) {
    kc, mods := ev.Detail, ev.State
    for _, m := range xgbutil.IgnoreMods {
        mods &= ^m
    }

    xu.RunKeyBindCallbacks(ev, xevent.KeyPress, ev.Event, mods, kc)
}

// RunKeyReleaseCallbacks infers the window, keycode and modifiers from a
// KeyPressEvent and runs the corresponding callbacks.
func RunKeyReleaseCallbacks(xu *xgbutil.XUtil, ev xevent.KeyReleaseEvent) {
    kc, mods := ev.Detail, ev.State
    for _, m := range xgbutil.IgnoreMods {
        mods &= ^m
    }

    xu.RunKeyBindCallbacks(ev, xevent.KeyRelease, ev.Event, mods, kc)
}
