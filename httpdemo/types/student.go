package types

import "net/http"

type Person struct{}

func (p *Person) Hello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello"))
}
