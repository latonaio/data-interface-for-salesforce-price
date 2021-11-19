package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	models "github.com/latonaio/salesforce-data-models"
	"github.com/latonaio/aion-core/pkg/go-client/msclient"
	"github.com/latonaio/data-interface-for-salesforce-price/internal/handlers"
	"github.com/latonaio/data-interface-for-salesforce-price/internal/resources"

	"github.com/latonaio/aion-core/pkg/log"
)

func main() {
	// Create Kanban client
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := models.NewDBConPool(ctx); err != nil {
		panic(err)
	}
	if err := newKanbanClient(ctx); err != nil {
		log.Fatalf("failed to get kanban client: %v", err)
	}
	log.Printf("successful get kanban client")
	defer kanbanClient.Close()

	// Get Kanban channel by Kanban client
	kanbanCh := kanbanClient.GetKanbanCh()
	log.Printf("successful get kanban channel")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM)
	for {
		select {
		case s := <-signalCh:
			fmt.Printf("received signal: %s", s.String())
			goto END
		case k := <-kanbanCh:
			if k == nil {
				goto NODATA
			}

			// Get metadata from Kanban
			fromMetadata, err := msclient.GetMetadataByMap(k)
			if err != nil {
				log.Printf("failed to get fromMetadata: %v", err)
				continue
			}

			log.Printf("got metadata from kanban")
			log.Printf("metadata: %v", fromMetadata)

			toMetadata, err := handle(fromMetadata)
			if err != nil {
				log.Printf("failed to handle: %v", err)
				continue
			}

			if toMetadata != nil {
				if err := writeKanban(toMetadata); err != nil {
					log.Printf("failed to write kanban: %v", err)
					continue
				}
				log.Printf("write metadata to kanban")
				log.Printf("metadata: %v", toMetadata)
			}
		}
	NODATA:
	}
END:
}

func handle(fromMetadata map[string]interface{}) (map[string]interface{}, error) {
	ck, ok := fromMetadata["connection_type"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid connection key")
	}

	if ck == "response" {
		key, ok := fromMetadata["key"].(string)
		if !ok {
			return nil, errors.New("invalid key")
		}
		if key == "PriceRecord" {
			if err := handlers.HandlePriceRecord(fromMetadata, writeKanban); err != nil {
				return nil, err
			}
		} else if key == "PriceRecordSeriesNumber" {
			if err := handlers.HandlePriceRecordSeriesNumber(fromMetadata); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	if ck == "request" {
		// Get Farm from metadata
		resource, err := resources.NewPriceRecord(fromMetadata)
		if err != nil {
			return nil, fmt.Errorf("failed to construct resource: %v", err)
		}

		// Build metadata for Kanban
		toMetadata, err := resource.BuildMetadata()
		if err != nil {
			return nil, fmt.Errorf("failed to build metadata: %v", err)
		}
		return toMetadata, nil
	}
	return nil, errors.New("failed to handle")
}
