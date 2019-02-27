// Copyright (c) 2019 KIDTSUNAMI
// Author: alex@kidtsunami.com
//

package blockwatch

import (
	"fmt"
	"time"
)

type CollapseMode string

const (
	CollapseInvalid        CollapseMode = ""
	CollapseNone           CollapseMode = "none"
	CollapseOneMinute      CollapseMode = "1m"
	CollapseFiveMinutes    CollapseMode = "5m"
	CollapseFifteenMinutes CollapseMode = "15m"
	CollapseThirtyMinutes  CollapseMode = "30m"
	CollapseOneHour        CollapseMode = "1h"
	CollapseThreeHours     CollapseMode = "3h"
	CollapseSixHours       CollapseMode = "6h"
	CollapseTwelveHours    CollapseMode = "12h"
	CollapseDaily          CollapseMode = "1d"
	CollapseWeekly         CollapseMode = "1w"
	CollapseMonthly        CollapseMode = "1M"
	CollapseQuarterly      CollapseMode = "3M"
	CollapseAnnual         CollapseMode = "1y"
)

func (m CollapseMode) String() string {
	return string(m)
}

func ParseCollapseModeIgnoreError(s string) CollapseMode {
	m, _ := ParseCollapseMode(s)
	return m
}

func ParseCollapseMode(s string) (CollapseMode, error) {
	switch s {
	case "oneminute", "1m":
		return CollapseOneMinute, nil
	case "fiveminutes", "5m":
		return CollapseFiveMinutes, nil
	case "fifteenminutes", "15m":
		return CollapseFifteenMinutes, nil
	case "thirtyminutes", "30m":
		return CollapseThirtyMinutes, nil
	case "onehour", "1h":
		return CollapseOneHour, nil
	case "threehours", "3h":
		return CollapseThreeHours, nil
	case "sixhours", "6h":
		return CollapseSixHours, nil
	case "twelvehours", "12h":
		return CollapseTwelveHours, nil
	case "daily", "1d":
		return CollapseDaily, nil
	case "weekly", "1w":
		return CollapseWeekly, nil
	case "monthly", "1M":
		return CollapseMonthly, nil
	case "quarterly", "1q", "3M":
		return CollapseQuarterly, nil
	case "annual", "1y":
		return CollapseAnnual, nil
	case "none", "":
		return CollapseNone, nil
	default:
		return CollapseInvalid, fmt.Errorf("invalid collapse mode '%s'", s)
	}
}

func NewCollapseMode(d time.Duration) CollapseMode {
	switch true {
	case d == 0:
		return CollapseNone
	case d <= time.Minute:
		return CollapseOneMinute
	case d <= 5*time.Minute:
		return CollapseFiveMinutes
	case d <= 15*time.Minute:
		return CollapseFifteenMinutes
	case d <= 30*time.Minute:
		return CollapseThirtyMinutes
	case d <= time.Hour:
		return CollapseOneHour
	case d <= 3*time.Hour:
		return CollapseThreeHours
	case d <= 6*time.Hour:
		return CollapseSixHours
	case d <= 12*time.Hour:
		return CollapseTwelveHours
	case d <= 24*time.Hour:
		return CollapseDaily
	case d <= 7*24*time.Hour:
		return CollapseWeekly
	case d <= 31*24*time.Hour: // approx, month is avg 30.41 days
		return CollapseMonthly
	case d <= 92*24*time.Hour: // approx, quarter is avg 91.25 days
		return CollapseQuarterly
	default:
		return CollapseAnnual
	}
}

func (m CollapseMode) IsValid() bool {
	return m != CollapseInvalid
}

// Text/JSON conversion
func (m CollapseMode) MarshalText() ([]byte, error) {
	return []byte(m), nil
}

func (m *CollapseMode) UnmarshalText(data []byte) error {
	mm, err := ParseCollapseMode(string(data))
	if err != nil {
		return err
	}
	*m = mm
	return nil
}

func (m CollapseMode) Duration() time.Duration {
	switch m {
	case CollapseOneMinute:
		return time.Minute // "1m"
	case CollapseFiveMinutes:
		return 5 * time.Minute // "5m"
	case CollapseFifteenMinutes:
		return 15 * time.Minute // "15m"
	case CollapseThirtyMinutes:
		return 30 * time.Minute // "30m"
	case CollapseOneHour:
		return time.Hour // "1h"
	case CollapseThreeHours:
		return 3 * time.Hour // "3h"
	case CollapseSixHours:
		return 6 * time.Hour // "6h"
	case CollapseTwelveHours:
		return 12 * time.Hour // "12h"
	case CollapseDaily:
		return 24 * time.Hour // "1d"
	case CollapseWeekly:
		return 24 * 7 * time.Hour // "1w"
	case CollapseMonthly:
		return 30 * 24 * time.Hour // "1M" // approx, month is avg 30.41 days
	case CollapseQuarterly:
		return 90 * 24 * time.Hour // "3M" // approx, quarter is avg 91.25 days
	case CollapseAnnual:
		return 365 * 24 * time.Hour // "1y"
	default:
		return time.Minute // default
	}
}
