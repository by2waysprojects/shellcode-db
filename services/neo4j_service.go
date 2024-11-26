package services

import (
	"context"
	"log"
	"shellcode-db/model"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jService struct {
	Driver neo4j.DriverWithContext
}

func NewNeo4jService(uri, username, password string) *Neo4jService {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatalf("Failed to create Neo4j driver: %v", err)
	}
	return &Neo4jService{Driver: driver}
}

func (s *Neo4jService) Close() {
	s.Driver.Close(context.Background())
}

func (s *Neo4jService) GetAllArchitectures() ([]model.Architecture, error) {
	ctx := context.Background()
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	architectures := []model.Architecture{}
	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, "MATCH (p:Architecture) RETURN p.id AS id, p.name AS name", nil)
		if err != nil {
			return nil, err
		}

		for result.Next(ctx) {
			record := result.Record()
			architecture := model.Architecture{
				ID:   record.Values[0].(string),
				Name: record.Values[1].(string),
			}
			architectures = append(architectures, architecture)
		}
		return nil, result.Err()
	})

	return architectures, err
}

func (s *Neo4jService) GetShellcodesByArchitectureID(architectureID string) ([]model.Shellcode, error) {
	ctx := context.Background()
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	shellcodes := []model.Shellcode{}
	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH (p:Architecture {id: $architectureID})-[:HAS_CHILD]->(c:Shellcode)
			RETURN c.id AS id, c.name AS name,  c.datePublished as datePublished, c.data as data
		`
		params := map[string]interface{}{"architectureID": architectureID}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		for result.Next(ctx) {
			record := result.Record()
			shellcode := model.Shellcode{
				ID:            record.Values[0].(string),
				Name:          record.Values[1].(string),
				DatePublished: record.Values[2].(string),
				Data:          record.Values[3].(string),
			}
			shellcodes = append(shellcodes, shellcode)
		}
		return nil, result.Err()
	})

	return shellcodes, err
}

func (s *Neo4jService) CreateArchitecture(architecture model.Architecture) error {
	ctx := context.Background()
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MERGE (p:Architecture {id: $id})
			SET p.name = $name
		`
		params := map[string]interface{}{"id": architecture.ID, "name": architecture.Name}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

func (s *Neo4jService) CreateShellcodeWithArchitecture(architectureID string, shellcode model.Shellcode) error {
	ctx := context.Background()
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MERGE (c:Shellcode {id: $shellcodeID})
			SET c.name = $shellcodeName, c.datePublished = $datePublished, c.data = $data
			WITH c
			MATCH (p:Architecture {id: $architectureID})
			MERGE (p)-[:HAS_CHILD]->(c)
		`
		params := map[string]interface{}{
			"architectureID": architectureID,
			"shellcodeID":    shellcode.ID,
			"shellcodeName":  shellcode.Name,
			"datePublished":  shellcode.DatePublished,
			"data":           shellcode.Data,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}
