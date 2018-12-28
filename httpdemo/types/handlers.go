package types

import (
	"fmt"
	"net/http"
)

func Say(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("斩断三千烦恼丝"))
}

func X(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(fmt.Sprintf("X = %d\n", -1)))
}