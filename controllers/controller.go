package controllers

import (
	"encoding/json"
	"net/http"
	"shellcode-db/services"

	"github.com/gorilla/mux"
)

type ShellCodeController struct {
	DBService         *services.Neo4jService
	DataLoaderService *services.DataLoaderService
}

func NewShellCodeController(dbService *services.Neo4jService, dataLoaderService *services.DataLoaderService) *ShellCodeController {
	return &ShellCodeController{DBService: dbService, DataLoaderService: dataLoaderService}
}

func (sc *ShellCodeController) GetAllArchitectures(w http.ResponseWriter, r *http.Request) {
	archiectures, err := sc.DBService.GetAllArchitectures()
	if err != nil {
		http.Error(w, "Failed to fetch architectures", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(archiectures)
}

func (sc *ShellCodeController) GetShellcodesByArchitectureID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	architectureID := vars["id"]

	shellcodes, err := sc.DBService.GetShellcodesByArchitectureID(architectureID)
	if err != nil {
		http.Error(w, "Failed to fetch shellcodes", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(shellcodes)
}

func (sc *ShellCodeController) LoadData(w http.ResponseWriter, r *http.Request) {
	err := sc.DataLoaderService.LoadData()
	if err != nil {
		http.Error(w, "Failed to load data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data successfully loaded"))
}
