package types

import "time"

type SessionsResponse struct {
	Cursor   uint64        `json:"cursor" example:"0"`
	Sessions []SessionInfo `json:"sessions"`
}

type SessionInfo struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	IP        string    `json:"ip"`
	Device    string    `json:"device"`
	LoginTime time.Time `json:"login_time"`
}

func (s SessionInfo) GetIP() string {
	return s.IP
}

func (s SessionInfo) GetDevice() string {
	return s.Device
}
