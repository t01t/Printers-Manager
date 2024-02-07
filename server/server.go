package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alexbrainman/printer"
	"github.com/gorilla/mux"
)

type Printer struct {
	Name    string
	Printer *printer.Printer
}

var DefaultPrinter *Printer

func Init() error {
	r := mux.NewRouter()
	r.HandleFunc("/printers", PrintersList)
	r.HandleFunc("/printers/{name}/setDefault", SetDefaultPrinter)
	r.HandleFunc("/printers/{name}/jobs", GetPrinterJobs)
	r.HandleFunc("/printer", GetDefaultPrinter)
	r.HandleFunc("/printer/jobs", GetDefaultPrinterJobs)
	r.HandleFunc("/printer/print", Print)

	return http.ListenAndServe("127.0.0.1:6969", r)
}

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

func GetDefaultPrinter(w http.ResponseWriter, r *http.Request) {
	if DefaultPrinter == nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("no default printer is selected")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, DefaultPrinter.Name)
}

func GetDefaultPrinterJobs(w http.ResponseWriter, r *http.Request) {
	if DefaultPrinter == nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("no default printer is selected")
		return
	}
	jobs, err := DefaultPrinter.Printer.Jobs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed getting '%s' jobs, %s \n", DefaultPrinter.Name, err.Error())
		return
	}

	jobsJson, err := json.Marshal(jobs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed exporting '%s' jobs to JSON, %s \n", DefaultPrinter.Name, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jobsJson))
}

func ClearDefaultPrinter(w http.ResponseWriter, r *http.Request) {
	if DefaultPrinter == nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
		return
	}
	DefaultPrinter.Printer.Close()
	DefaultPrinter = nil

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func SetDefaultPrinter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	p, err := printer.Open(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to open printer,", err.Error())
		return
	}
	DefaultPrinter = &Printer{
		Name:    name,
		Printer: p,
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func Print(w http.ResponseWriter, r *http.Request) {
	paths := []string{}

	err := json.NewDecoder(r.Body).Decode(&paths)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to parse request body, ", err.Error())
		return
	}
	fmt.Println(paths)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
