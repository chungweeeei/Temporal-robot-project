package activities

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
)

func (ra *RobotActivities) SimpleAction(ctx context.Context) (string, error) {

	logger := activity.GetLogger(ctx)
	logger.Info("Standup Activity started. Waiting/Blocking...")

	// 使用 Ticker 每秒送一次心跳，確保 Server 知道我們還活著
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// 總體超時控制器 (例如 10 分鐘後自動結束)
	timeout := time.After(1 * time.Minute)

	for {
		select {
		case <-ticker.C:
			// 每秒觸發一次，送心跳
			logger.Info("Simple Action send heartbeat")
			activity.RecordHeartbeat(ctx, "waiting")
			// 這裡不 return，繼續迴圈
			if ctx.Err() != nil {
				logger.Info("Detected context error right after heartbeat", "error", ctx.Err())
				return "", ctx.Err()
			}
		case <-timeout:
			// 10 分鐘到了，正常結束
			logger.Info("Simple Action finished normally (time up)")
			return "finished", nil

		case <-ctx.Done():
			// 收到 Cancel 信號
			logger.Info("Simple Action has been cancelled")
			return "", ctx.Err()
		}
	}
}
