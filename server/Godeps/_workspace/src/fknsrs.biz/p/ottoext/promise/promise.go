package promise

import (
	"github.com/GeertJohan/go.rice"
	"github.com/robertkrimen/otto"

	"fknsrs.biz/p/ottoext/loop"
	"fknsrs.biz/p/ottoext/timers"
)

func Define(vm *otto.Otto, l *loop.Loop) error {
	if v, err := vm.Get("Promise"); err != nil {
		return err
	} else if !v.IsUndefined() {
		return nil
	}

	if err := timers.Define(vm, l); err != nil {
		return err
	}

	s, err := vm.Compile("bundle.js", rice.MustFindBox("dist-promise").MustString("bundle.js"))
	if err != nil {
		return err
	}

	if _, err := vm.Run(s); err != nil {
		return err
	}

	return nil
}
