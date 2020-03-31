package ssdp

import (
	"context"
	"time"

	log "github.com/JK-97/edge-guard/lowapi/logger"

	"github.com/koron/go-ssdp"
)

type ssdpClient struct {
	WorkerID string
	Interval int
}

func (c *ssdpClient) Aliving(ctx context.Context) error {

	ad, err := ssdp.Advertise(
		"EDGE-GUARD:EDGE-GUARD", // send as "ST"
		"alive:"+c.WorkerID,     // send as "USN"
		"",                      // send as "LOCATION"
		"",                      // send as "SERVER"
		1800)                    // send as "maxAge" in "CACHE-CONTROL"
	if err != nil {
		return err
	}
	aliveTicker := time.NewTicker(time.Duration(c.Interval) * time.Second)
	defer aliveTicker.Stop()

	for {
		select {
		case <-aliveTicker.C:
			err := ad.Alive()
			if err != nil {
				log.Info(err)
			}
		case <-ctx.Done():
			return nil

		}
	}

}

func NewClient(workerID string, interval int) *ssdpClient {
	return &ssdpClient{
		WorkerID: workerID,
		Interval: interval,
	}
}
