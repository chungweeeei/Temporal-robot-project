package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/activities"
	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/helper"
	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/workflows"
	"github.com/chungweeeei/Temporal-robot-project/pkg"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {

	// Register low-level slog hanlder setting level to INFO
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Register Temporal client
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
		Logger:   logger,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// check robot ip
	robotIP := os.Getenv("ROBOT_IP")
	if robotIP == "" {
		robotIP = "localhost"
	}
	// create StatusCache instance
	statusCache := &activities.StatusCache{}

	// Background subscriber
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go RobotStatusSubscriber(ctx, fmt.Sprintf("ws://%s:9090/", robotIP), statusCache)

	w := worker.New(c, "ROBOT_TASK_QUEUE", worker.Options{})

	activities := activities.NewRobotActivities(robotIP, statusCache)
	w.RegisterWorkflow(workflows.RobotWorkflow)
	w.RegisterActivity(activities)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

func RobotStatusSubscriber(
	ctx context.Context,
	wsURL string,
	cache *activities.StatusCache,
) {

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := subscribeLoop(ctx, wsURL, cache); err != nil {
					log.Println("Subscriber error, reconnecting in 5s:", err)
					time.Sleep(5 * time.Second)
				}
			}
		}
	}()
}

// 定義一個能夠容忍型別混亂的結構 (應該要從 ROS端直接修改這樣上層就不用多做一層parsing)
type RawRobotStatus struct {
	ApiID        int         `json:"api_id"`
	BatteryLevel interface{} `json:"battery_level"` // Could be string "94" or int 94
	Pose         struct {
		Position struct {
			X interface{} `json:"x"` // string or number
			Y interface{} `json:"y"`
			Z interface{} `json:"z"`
		} `json:"position"`
		Orientation struct {
			X interface{} `json:"x"`
			Y interface{} `json:"y"`
			Z interface{} `json:"z"`
			W interface{} `json:"w"`
		} `json:"orientation"`
	} `json:"pose"`
	MissionID interface{} `json:"mission_id"`
	Mission   struct {
		Code    interface{} `json:"code"`    // string or int
		Message interface{} `json:"message"` // string
	} `json:"mission"`
}

func subscribeLoop(
	ctx context.Context,
	wsURL string,
	cache *activities.StatusCache,
) error {

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 1. 發送訂閱請求
	subscribeMsg := map[string]interface{}{
		"op":            "subscribe",
		"topic":         "/api/info",
		"type":          "std_msgs/msg/String",
		"throttle_rate": 0,
		"queue_length":  1,
	}
	if err := conn.WriteJSON(subscribeMsg); err != nil {
		return err
	}

	// 2. 持續接收訊息
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				return err
			}

			var msg pkg.TopicResponse
			if err := json.Unmarshal(message, &msg); err != nil {
				continue
			}

			var resp struct {
				DeviceName   string         `json:"device_name"`
				DeviceStatus RawRobotStatus `json:"device_status"`
				TimeStamp    string         `json:"timestamp"`
			}
			if err := json.Unmarshal([]byte(msg.Msg.Data), &resp); err != nil {
				continue
			}

			if resp.DeviceStatus.BatteryLevel == nil {
				continue
			}

			status := activities.RobotStatus{
				ApiID:        resp.DeviceStatus.ApiID,
				BatteryLevel: helper.ToInt(resp.DeviceStatus.BatteryLevel),
			}
			status.Pose.Position.X = helper.ToFloat(resp.DeviceStatus.Pose.Position.X)
			status.Pose.Position.Y = helper.ToFloat(resp.DeviceStatus.Pose.Position.Y)
			status.Pose.Position.Z = helper.ToFloat(resp.DeviceStatus.Pose.Position.Z)

			status.Pose.Orientation.X = helper.ToFloat(resp.DeviceStatus.Pose.Orientation.X)
			status.Pose.Orientation.Y = helper.ToFloat(resp.DeviceStatus.Pose.Orientation.Y)
			status.Pose.Orientation.Z = helper.ToFloat(resp.DeviceStatus.Pose.Orientation.Z)
			status.Pose.Orientation.W = helper.ToFloat(resp.DeviceStatus.Pose.Orientation.W)

			status.MissionID = resp.DeviceStatus.MissionID.(string)
			status.Mission.Code = activities.MissionCode(helper.ToInt(resp.DeviceStatus.Mission.Code))
			status.Mission.Message = resp.DeviceStatus.Mission.Message.(string)

			color.Green("Robot current MissionID: %s, Code: %d, Message: %s", status.MissionID, status.Mission.Code, status.Mission.Message)

			// update cache value
			cache.Update(status)
		}
	}
}
