package routes

import (
	"shellcode-db/controllers"

	"github.com/gorilla/mux"
)

func SetupRoutes(
	shellCodeController *controllers.ShellCodeController,
) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/architectures", shellCodeController.GetAllArchitectures).Methods("GET")
	router.HandleFunc("/architectures/{id}/shellcodes", shellCodeController.GetShellcodesByArchitectureID).Methods("GET")
	router.HandleFunc("/load-data", shellCodeController.LoadData).Methods("POST")

	return router
}
