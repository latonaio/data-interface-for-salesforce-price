package resources

import (
	"errors"
	"fmt"
)

type PriceRecordSeriesNumber struct {
	method   string
	metadata map[string]interface{}
}

func (c *PriceRecordSeriesNumber) objectName() string {
	const obName = "PriceRecordSeriesNumber"
	return obName
}

func NewPriceRecordSeriesNumber(metadata map[string]interface{}) (*PriceRecordSeriesNumber, error) {
	rawMethod, ok := metadata["method"]
	if !ok {
		return nil, errors.New("missing required parameters: method")
	}
	method, ok := rawMethod.(string)
	if !ok {
		return nil, errors.New("failed to convert interface{} to string")
	}
	return &PriceRecordSeriesNumber{
		method:   method,
		metadata: metadata,
	}, nil
}

func (c *PriceRecordSeriesNumber) getMetadata() (map[string]interface{}, error) {
	prId := c.metadata["priceRecordId"].(string)
	query := map[string]string{
		"priceRecordId": prId,
	}
	return buildMetadata(c.method, c.objectName(), priceConnectionKey, "", query, ""), nil
}

func (c *PriceRecordSeriesNumber) BuildMetadata() (map[string]interface{}, error) {
	switch c.method {
	case "get":
		return c.getMetadata()
	}
	return nil, fmt.Errorf("invalid method: %s", c.method)
}
