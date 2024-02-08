package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/alexbrainman/printer"
	"github.com/gorilla/mux"
	"github.com/t01t/printers-manager/print"
)

func PrintersList(w http.ResponseWriter, r *http.Request) {
	printers, err := printer.ReadNames()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to get printers,", err.Error())
		return
	}
	printersJson, err := json.Marshal(printers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed parse printers list,", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(printersJson))
}

func GetPrinterJobs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	p, err := printer.Open(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to open '%s', %s \n", name, err.Error())
		return
	}
	defer p.Close()

	jobs, err := p.Jobs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed getting '%s' jobs, %s \n", name, err.Error())
		return
	}

	jobsJson, err := json.Marshal(jobs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed exporting '%s' jobs to JSON, %s \n", name, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jobsJson))
}

func PrintFromPaths(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	path, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to parse request body, ", err.Error())
		return
	}
	fmt.Println(string(path))

	err = print.AddFileByPath(name, string(path))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to add to jobs '%s', %s\n", string(path), err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
