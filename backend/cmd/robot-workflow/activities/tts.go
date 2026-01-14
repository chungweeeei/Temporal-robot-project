package activities

import (
	"context"
	"encoding/json"
	"fmt"
)

func (ra *RobotActivities) TTS(ctx context.Context, params map[string]interface{}) (string, error) {

	text, ok := params["text"].(string)
	if !ok {
		return "", fmt.Errorf("invalid parameters for TTS activity")
	}

	return executeWithHeartbeat(ctx, func() (string, error) {

		data := map[string]interface{}{
			"api_id":     RobotTTSCommandID,
			"text":       text,
			"voice_name": "Mandarin-CN.Male-1",
			"speed":      1.0,
			"volume":     1.0,
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			return "", err
		}

		return ra.Client.CallService(ctx, "TTS", string(dataBytes))
	})
}
