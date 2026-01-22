package activities

import (
	"errors"
	"sync"
	"time"

	config "github.com/chungweeeei/Temporal-robot-project/internal/config/activity"
)

var (
	ErrStatusNotAvailable = errors.New("robot status not available yet")
	ErrStatusStale        = errors.New("robot status is stale")
)

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

type RobotStatus struct {
	ApiID        int `json:"api_id"`
	BatteryLevel int `json:"battery_level"`
	Pose         struct {
		Orientation struct {
			W float64 `json:"w"`
			X float64 `json:"x"`
			Y float64 `json:"y"`
			Z float64 `json:"z"`
		} `json:"orientation"`
		Position struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
			Z float64 `json:"z"`
		} `json:"position"`
	} `json:"pose"`
	MissionID string `json:"mission_id"`
	Mission   struct {
		Code    config.MissionCode `json:"code"`
		Message string             `json:"message"`
	} `json:"mission"`
}

// Status Cache (Background goroutine updates this periodically)
type CacheStatus struct {
	mu          sync.RWMutex
	status      RobotStatus
	lastUpdated time.Time
	initialized bool
}

func NewCacheStatus() *CacheStatus {
	return &CacheStatus{}
}

func (c *CacheStatus) Get() (RobotStatus, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.initialized {
		return RobotStatus{}, ErrStatusNotAvailable
	}

	if time.Since(c.lastUpdated) > 10*time.Second {
		return c.status, ErrStatusStale
	}

	return c.status, nil
}

func (c *CacheStatus) Update(status RobotStatus) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.status = status
	c.lastUpdated = time.Now()
	c.initialized = true
}

func (c *CacheStatus) IsReady() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.initialized
}
