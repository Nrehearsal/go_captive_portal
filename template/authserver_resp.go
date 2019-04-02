package template

import "time"

type OnlineUser struct {
	Username  string    `json:"username"`
	Level     int       `json:"level"`
	IP        string    `json:"ip"`
	Mac       string    `json:"mac"`
	ExpiredAt time.Time `json:"expired_at"`
}
