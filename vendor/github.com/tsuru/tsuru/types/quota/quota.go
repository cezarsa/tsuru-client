// Copyright 2018 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quota

import (
	"errors"
	"fmt"
)

type Quota struct {
	Limit int `json:"limit"`
	InUse int `json:"inuse"`
}

// UnlimitedQuota is the struct which any new unlimited quota copies from.
var UnlimitedQuota = Quota{Limit: -1, InUse: 0}

func (q *Quota) IsUnlimited() bool {
	return -1 == q.Limit
}

type QuotaService interface {
	Inc(name string, delta int) error
	Set(name string, quantity int) error
	SetLimit(name string, limit int) error
	Get(name string) (*Quota, error)
}

type QuotaStorage interface {
	Inc(name string, delta int) error
	SetLimit(name string, limit int) error
	Get(name string) (*Quota, error)
	Set(name string, quantity int) error
}

type QuotaExceededError struct {
	Requested uint
	Available uint
}

func (err *QuotaExceededError) Error() string {
	return fmt.Sprintf("Quota exceeded. Available: %d, Requested: %d.", err.Available, err.Requested)
}

var (
	ErrNotEnoughReserved       = errors.New("Not enough reserved items")
	ErrLimitLowerThanAllocated = errors.New("New limit is less than the current allocated value")
	ErrLessThanZero            = errors.New("Invalid value, cannot be less than 0")
	ErrQuotaNotFound           = errors.New("quota not found")
)
