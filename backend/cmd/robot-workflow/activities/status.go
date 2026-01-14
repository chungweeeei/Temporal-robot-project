package activities

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrStatusNotAvailable = errors.New("robot status not available yet")
	ErrStatusStale        = errors.New("robot status is stale")
)

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
}

// Shared Status Cache (Backgroud goroutine updates this cache periodically)
type StatusCache struct {
	mu          sync.RWMutex
	status      RobotStatus
	lastUpdated time.Time
	initialized bool
}

func NewStatusCache() *StatusCache {
	return &StatusCache{}
}

func (c *StatusCache) Get() (RobotStatus, error) {
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

func (c *StatusCache) Update(s RobotStatus) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status = s
	c.lastUpdated = time.Now()
	c.initialized = true
}

func (c *StatusCache) IsReady() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.initialized
}

func (ra *RobotActivities) GetStatus(ctx context.Context) (RobotStatus, error) {

	if ra.StatusCache == nil {
		return RobotStatus{}, ErrStatusNotAvailable
	}

	return ra.StatusCache.Get()
}
