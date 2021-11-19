package resources

import (
	"errors"
	"fmt"
)

type PriceRecord struct {
	method   string
	metadata map[string]interface{}
}

func (c *PriceRecord) objectName() string {
	const obName = "PriceRecord"
	return obName
}

func NewPriceRecord(metadata map[string]interface{}) (*PriceRecord, error) {
	rawMethod, ok := metadata["method"]
	if !ok {
		return nil, errors.New("missing required parameters: method")
	}
	method, ok := rawMethod.(string)
	if !ok {
		return nil, errors.New("failed to convert interface{} to string")
	}
	return &PriceRecord{
		method:   method,
		metadata: metadata,
	}, nil
}

func (c *PriceRecord) getMetadata() (map[string]interface{}, error) {
	dId := c.metadata["districtId"].(string)
	query := map[string]string{
		"platId": dId,
	}
	return buildMetadata(c.method, c.objectName(), priceConnectionKey, "", query, ""), nil
}

func (c *PriceRecord) BuildMetadata() (map[string]interface{}, error) {
	switch c.method {
	case "get":
		return c.getMetadata()
	}
	return nil, fmt.Errorf("invalid method: %s", c.method)
}
