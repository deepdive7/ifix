package main

import (
	"fmt"
	"github.com/deepdive7/ifix"
	"github.com/deepdive7/ifix/httpdemo/types"
	"net/http"
	"reflect"
)

func Info(info map[string]string) {
	fmt.Println("Loading...")
	info["PatchPerson"] = "PatchPerson"
	info["PatchSay"] = "PatchSay"
	info["PatchX"] = "PatchX"
}

func Say(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("木鱼孤灯清白衣"))
}

func PatchPerson(a *types.Person) {
	rt := reflect.TypeOf(a)
	ifix.PatchInstanceMethod(rt, "Hello", func(p *types.Person, w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("patch2 running"))
	})
}

func PatchSay(h func(w http.ResponseWriter, req *http.Request)) {
	ifix.Patch(h, Say)
}

func PatchX(handler func(w http.ResponseWriter, req *http.Request)) {
	ifix.Patch(handler, func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(fmt.Sprintf("X = %d\n", 520)))
	})
}