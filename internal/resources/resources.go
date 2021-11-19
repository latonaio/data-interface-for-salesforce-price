package resources

type Resource interface {
	getMetadata() (map[string]interface{}, error)
	postMetadata() (map[string]interface{}, error)
	updateMetadata() (map[string]interface{}, error)
	BuildMetadata() (map[string]interface{}, error)
	BuildConnectionKey() string
}

const priceConnectionKey = "price_get"

func buildMetadata(method, object, connectionKey, pathParam string, queryParams map[string]string, body string) map[string]interface{} {
	metadata := map[string]interface{}{
		"method":         method,
		"object":         object,
		"connection_key": connectionKey,
	}
	if len(pathParam) > 0 {
		metadata["path_param"] = pathParam
	}
	if queryParams != nil {
		metadata["query_params"] = queryParams
	}
	if body != "" {
		metadata["body"] = body
	}
	return metadata
}
