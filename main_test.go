package main

import (
	"reflect"
	"testing"
)

func TestMainProcessing(t *testing.T) {

	tests := []struct {
		FromMetadata     map[string]interface{}
		ExpectToMetadata map[string]interface{}
	}{
		// PriceRecord
		{
			FromMetadata: map[string]interface{}{
				"connection_type": "request",
				"districtId":      "a080l000007iaF4AAI",
				"method":          "get",
				"object":          "Price",
			},
			ExpectToMetadata: map[string]interface{}{
				"connection_key": "price_get",
				"method":         "get",
				"object":         "PriceRecord",
				"query_params": map[string]string{
					"platId": "a080l000007iaF4AAI",
				},
			},
		},
	}
	for i, tt := range tests {
		got, err := handle(tt.FromMetadata)
		if err != nil {
			t.Errorf("faild to handle: %v", err)
		}

		if !reflect.DeepEqual(got, tt.ExpectToMetadata) {
			t.Errorf("%d# metadata got = %v, want %v", i, got, tt.ExpectToMetadata)
		}
	}
}
