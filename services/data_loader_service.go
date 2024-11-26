package services

import (
	"log"
)

type DataLoaderService struct {
	APIService   *APIService
	Neo4jService *Neo4jService
}

func NewDataLoaderService(apiService *APIService, neo4jService *Neo4jService) *DataLoaderService {
	return &DataLoaderService{
		APIService:   apiService,
		Neo4jService: neo4jService,
	}
}

func (dl *DataLoaderService) LoadData() error {
	// Fetch Architectures from API
	architectures, err := dl.APIService.FetchArchitectures()
	if err != nil {
		return err
	}

	// Insert Architectures into Neo4j
	for _, architecture := range architectures {
		err := dl.Neo4jService.CreateArchitecture(architecture)
		if err != nil {
			log.Printf("Failed to insert architecture %s: %v", architecture.ID, err)
			continue
		}

		for _, shellcode := range architecture.Shellcodes {
			if err := dl.Neo4jService.CreateShellcodeWithArchitecture(architecture.ID, shellcode); err != nil {
				log.Printf("Failed to insert shellcode %s: %v", shellcode.ID, err)
			}
		}
	}
	return nil
}
