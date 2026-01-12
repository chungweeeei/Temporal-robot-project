package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chungweeeei/Temporal-robot-project/pkg"
	"go.temporal.io/sdk/activity"
)

const (
	ACTIVITY_HEARTBEAT_INTERVAL = 3 // unit in seconds
)

func generatePayload(data string) ([]byte, error) {
	req := pkg.ServiceRequest{
		Op:      "call_service",
		Service: "/api/system",
		Type:    "custom_msgs/srv/Api",
		Args: struct {
			Data string `json:"data"`
		}{
			Data: data,
		},
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func parseResponse(msg []byte) (string, error) {
	var resp pkg.ServiceResponse
	if err := json.Unmarshal(msg, &resp); err != nil {
		return "", err
	}
	if !resp.Result {
		return "", fmt.Errorf("server error")
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

/*
	Go原生select:
  		- 使用時機: 用於標準Go程式碼中，包括main func, goroutines以及Temporal Activities．
  		- Activity 本質上就是普通的 Go 函式，它運行在標準的 Go Runtime 上，不受 Temporal 的確定性 (Determinism) 限制。
	workflow.Selector:
  		- 使用時機: 只能用於 Temporal Workflow 定義中。
  		- 原因: Workflow 必須是確定性的 (Deterministic)。Go 原生的 select 在處理多個 channel 同時有資料時，其選取順序是隨機的 (Random).
*/

/*
錯誤範例：不能在 Workflow 裡這樣寫，因為 select 是隨機的
select {
case <-future1.Get(ctx, nil): // 這是 Go 原生寫法，在 Workflow 會出錯或不穩定
case <-future2.Get(ctx, nil):
}

// 正確範例：在 Workflow 裡要用 Selector
selector := workflow.NewSelector(ctx)
selector.AddFuture(future1, func(f workflow.Future) {
    // 處理 future1
})
selector.AddFuture(future2, func(f workflow.Future) {
    // 處理 future2
})
selector.Select()
*/
