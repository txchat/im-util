package rate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseDuration(str string) (time.Duration, error) {
	return time.ParseDuration(str)
}

func ParseRateString(str string) (int, time.Duration, error) {
	items := strings.Split(str, "/")
	if len(items) < 2 {
		return 0, 0, fmt.Errorf("cannot split rate need %d, got %d", 2, len(items))
	}
	num, err := strconv.ParseInt(items[0], 10, 32)
	if err != nil {
		return 0, 0, err
	}
	reg := regexp.MustCompile(`^[a-z]+$`)
	if !reg.MatchString(items[1]) {
		return 0, 0, fmt.Errorf("time not mach, got %s", items[1])
	}
	td, err := ParseDuration("1" + items[1])
	if err != nil {
		return 0, 0, err
	}
	return int(num), td, nil
}

func GetUserIntervalAndTimeInterval(userNum, msgNum int, tm time.Duration) (time.Duration, time.Duration) {
	userInterval := time.Duration(int64(tm) / int64(userNum))
	timeInterval := time.Duration(int64(tm) / int64(msgNum))
	return userInterval, timeInterval
}
