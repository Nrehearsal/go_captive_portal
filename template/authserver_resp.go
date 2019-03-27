package template

import "time"

type OnlineUse struct {
	Username  string    `json:"username"`
	IP        string    `json:"ip"`
	Mac       string    `json:"mac"`
	ExpiredAt time.Time `json:"expired_at"`
}
