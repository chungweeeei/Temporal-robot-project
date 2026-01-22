package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/activities"
	"github.com/chungweeeei/Temporal-robot-project/cmd/robot-workflow/workflows"
	"github.com/chungweeeei/Temporal-robot-project/pkg"
	"github.com/gorilla/websocket"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {

	// Register Temporal client
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Check robot ip settings in environment variables
	robotIP := os.Getenv("ROBOT_IP")
	if robotIP == "" {
		robotIP = "localhost"
	}

	// Register StatusCache instance
	statusCache := &activities.StatusCache{}

	// Background go routine for robot status subscription
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go RobotStatusSubscriber(ctx, fmt.Sprintf("ws://%s:9090/", robotIP), statusCache)

	// Register temporal worker
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
}

type RawRobotStatus struct {
	ApiID        int         `json:"api_id"`
	BatteryLevel interface{} `json:"battery_level"` // Could be string "94" or int 94
	Pose         struct {
		Position struct {
			X interface{} `json:"x"`
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
		Code    interface{} `json:"code"`
		Message interface{} `json:"message"`
	} `json:"mission"`
}

func subscribeLoop(
	ctx context.Context,
	wsURL string,
	cache *activities.CacheStatus,
) error {

	// Regsiter another websocket session
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 1. Publish subscribe message
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

	// 2. Continuously read message
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

			// Update cache value
			cache.Update(status)
		}
	}
}
