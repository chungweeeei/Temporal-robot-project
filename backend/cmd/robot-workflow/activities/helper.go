package activities

import (
	"encoding/json"
	"fmt"
)

func generateCommand(data string) ([]byte, error) {
	req := ServiceRequest{
		Op:      "call_service",
		Service: "/api/system",
		Type:    "custom_msgs/srv/Api",
		Args: RequestArgs{
			Data: data,
		},
	}

	msgBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return msgBytes, nil
}

func parseResponse(msg []byte) (string, error) {
	var resp ServiceResponse
	if err := json.Unmarshal(msg, &resp); err != nil {
		return "", err
	}
	if !resp.Result {
		return "", fmt.Errorf("server error")
	}
	return resp.Values.Data, nil
}
