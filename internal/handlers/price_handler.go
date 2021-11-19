package handlers

import (
	"errors"
	"fmt"

	"github.com/latonaio/data-interface-for-salesforce-price/internal/resources"
	"github.com/latonaio/salesforce-data-models"
)

func HandlePriceRecord(metadata map[string]interface{}, writeMetadata func(data map[string]interface{}) error) error {
	priceMasters, err := models.MetadataToPriceRecords(metadata)
	if err != nil {
		return fmt.Errorf("failed to convert metadata to models: %v", err)
	}
	if err := models.RegisterPriceRecordsAndCacheClear(priceMasters); err != nil {
		return errors.New("failed to cache clear and register table: " + err.Error())
	}
	if priceMasters[0].SfPriceRecordId == nil {
		return errors.New("price record id not found")
	}
	m := map[string]interface{}{
		"method":        "get",
		"priceRecordId": *priceMasters[0].SfPriceRecordId, //一件のみのはず
	}
	prn, err := resources.NewPriceRecordSeriesNumber(m)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to construct resource: %v", err))
	}
	toMetadata, err := prn.BuildMetadata()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to build metadata: %v", err))
	}
	if err := writeMetadata(toMetadata); err != nil {
		return errors.New(fmt.Sprintf("failed to write kanban: %v", err))
	}
	return nil
}

func HandlePriceRecordSeriesNumber(metadata map[string]interface{}) error {
	prsn, err := models.MetadataToPriceRecordSeriesNumbers(metadata)
	if err != nil {
		return fmt.Errorf("failed to convert metadata to models: %v", err)
	}
	if err := models.RegisterPriceRecordSeriesNumbersAndCacheClear(prsn); err != nil {
		return errors.New("failed to cache clear and register table: " + err.Error())
	}
	return nil
}
