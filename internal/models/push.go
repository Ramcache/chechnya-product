package models

import "github.com/SherClockHolmes/webpush-go"

type Subscription struct {
	Endpoint string
	P256dh   string
	Auth     string
	IsAdmin  bool
}

type PushSubscriptionRequest struct {
	Subscription webpush.Subscription `json:"subscription"`
	IsAdmin      bool                 `json:"is_admin"`
}

//test
