package activities

import (
	"context"
	"encoding/json"
	"time"

	"go.temporal.io/sdk/activity"
)

func (ra *RobotActivities) Sitdown(ctx context.Context, url string) (string, error) {

	logger := activity.GetLogger(ctx)

	errorChan := make(chan error)
	responseChan := make(chan string)

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	go func() {
		// execute business logic
		data := map[string]int{
			"api_id": RobotMotionControlID,
			"action": SitDownActionID,
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			errorChan <- err
			return
		}

		response, err := ra.Client.CallService(ctx, url, string(dataBytes))
		logger.Info("Receive Sitdown Service Response ", response)
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- response
	}()

	for {
		select {
		case <-ticker.C:
			// 定期發送心跳
			// 這是讓 Worker 有機會從 Server 收到 Cancel 指令的關鍵
			activity.RecordHeartbeat(ctx, "waiting-response")
		case err := <-errorChan:
			return "", err
		case response := <-responseChan:
			return response, nil
		case <-ctx.Done():
			logger.Info("Sitdown activity has been cancelled")
			return "", ctx.Err()
		}
	}
}
