package models

import "time"

type ObjectInfo struct {
	Name        string    `json:"name"`
	Key         string    `json:"key"`
	Size        int64     `json:"size"`
	Hash        string    `json:"hash"`
	RealPath    string    `json:"real_path"`
	ContextType string    `json:"context_type"`
	Version     string    `json:"version"`
	CreateAt    time.Time `json:"create_at"`
}
