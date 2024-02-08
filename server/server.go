package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Init() error {
	r := mux.NewRouter()
	r.HandleFunc("/printers", PrintersList)
	r.HandleFunc("/printers/{name}/jobs", GetPrinterJobs)
	r.HandleFunc("/printers/{name}/print", PrintFromPaths)

	return http.ListenAndServe("127.0.0.1:6969", r)
}
