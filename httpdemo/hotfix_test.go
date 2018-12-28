package hotfix

import (
	"fmt"
	"github.com/deepdive7/ifix"
	"github.com/deepdive7/ifix/httpdemo/types"
	"log"
	"net/http"
	"testing"
)

var s = &types.Person{}

func RegisterPerson() {
	http.HandleFunc("/hello", s.Hello)
}

func Patch(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	libPath := req.Form.Get("lib")
	//libPath := "./patch/patch.so"
	err := ifix.LoadDll(libPath, map[string][]interface{}{
		"PatchPerson": {&types.Person{}},
		"PatchSay":    {types.Say},
		"PatchX":      {types.X},
	})
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("Applied Patch and PatchSay"))
}

func TestPatch(t *testing.T) {
	RegisterPerson()
	http.HandleFunc("/say", types.Say)
	http.HandleFunc("/patch", Patch)
	http.HandleFunc("/x", types.X)
	log.Println("Server Running")
	fmt.Println(http.ListenAndServe(":8081", nil))
}
