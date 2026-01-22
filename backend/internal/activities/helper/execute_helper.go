package activities

import (
	"context"
	"encoding/json"
	"time"

	"github.com/chungweeeei/Temporal-robot-project/pkg"
	"go.temporal.io/sdk/activity"
)

const (
	ACTIVITY_HEARTBEAT_INTERVAL = 3
)

func generatePayload(actionType string, data any) ([]byte, error) {

	var req pkg.RobotServiceRequest
	switch actionType {
	case "Standup", "Sitdown", "Move", "TTS", "Status", "Stop":
		req.Op = "call_service"
		req.Service = "/api/system"
		req.Type = "custom_msgs/srv/Api"
		req.Args = struct {
			Data any `json:"data"`
		}{
			Data: data,
		}
	case "Head":
		req.Op = "call_service"
		req.Service = "/set_angle_tag"
		req.Type = "custom_msgs/srv/SetFloat"
		req.Args = struct {
			Data any `json:"data"`
		}{
			Data: data,
		}
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func parseResponse(msg []byte) (string, error) {
	var resp pkg.RobotServiceResponse
	if err := json.Unmarshal(msg, &resp); err != nil {
		return "", err
	}
	return resp.Values.Data, nil
}

// 整合這種自定義 schema，使用 Golang 的 Generics (泛型) 是最完美的解決方案。
// [Any] 表示這個函式接受任何型別 T
func executeWithHeartbeat[T any](ctx context.Context, operation func() (T, error)) (T, error) {

	logger := activity.GetLogger(ctx)

	// 使用 var zero T 來宣告該型別的零值，以便在錯誤時回傳
	var zero T

	resultCh := make(chan T, 1)
	errorCh := make(chan error, 1)

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	go func() {
		res, err := operation()
		if err != nil {
			errorCh <- err
			return
		}
		resultCh <- res
	}()

	for {
		select {
		case <-ticker.C:
			activity.RecordHeartbeat(ctx, "processing")
		case err := <-errorCh:
			return zero, err
		case res := <-resultCh:
			return res, nil
		case <-ctx.Done():
			logger.Info("Activity has been cancelled (ctx.Done)")
			return zero, ctx.Err()
		}
	}
}
