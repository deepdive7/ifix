package ifix

import (
	"fmt"
	"plugin"
	"reflect"
)

func Apply(a interface{}, b []interface{}) ([]reflect.Value){
	args := []reflect.Value{}
	for _, v := range b {
		t := v
		args = append(args, reflect.ValueOf(t))
	}
	return reflect.ValueOf(a).Call(args)
}

// 在同一个Plugin目录下编译的so不管名字还是存放路径都会表示为同一个plugin,只能加载一次，改变Plugin目录即可进行第二次编译，然后改变
// Case 1 (成功): 函数相同,分别是patch1/patch1.so 和patch2/patch2.so
// Case 2 (失败): 一次将两个库移动到patch/patch.so
func LoadDll(libPath string, patchers map[string][]interface{}) error {
	p, err := plugin.Open(libPath)
	if err != nil {
		println(err.Error())
		return err
	}
	info, err := p.Lookup("Info")
	patcherNames := map[string]string{}
	if err == nil {
		info.(func(map[string]string))(patcherNames)
		info = nil
		for k :=range patcherNames {
			fmt.Println("Patcher->", k)
		}
	}

	for name, args := range patchers{
		f, err := p.Lookup(patcherNames[name])
		if err != nil {
			println(err.Error())
			return err
		}
		fmt.Println("Apply ", name)
		Apply(interface{}(f), args)
	}
	p = nil
	return nil
}