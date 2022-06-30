package rate

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestParseRateString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name       string
		args       args
		want       int
		want1      time.Duration
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success",
			args: args{
				str: "600/m",
			},
			want:    600,
			want1:   time.Minute,
			wantErr: false,
		}, {
			name: "failed",
			args: args{
				str: "600/1m",
			},
			want:       0,
			want1:      0,
			wantErr:    true,
			wantErrMsg: fmt.Errorf("time not mach, got %s", "1m").Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseRateString(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRateString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if err.Error() != tt.wantErrMsg {
					t.Errorf("ParseRateString() error = %v, wantErrMsg %v", err, tt.wantErrMsg)
					return
				}
			}
			if got != tt.want {
				t.Errorf("ParseRateString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseRateString() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetUserIntervalAndTimeInterval(t *testing.T) {
	count := int32(0)
	userNum := 10000
	//msgNum := 600
	//tm := time.Second * 60

	ctx, closer := context.WithTimeout(context.Background(), time.Second*10)
	defer closer()

	wg := sync.WaitGroup{}
	for i := 0; i < userNum; i++ {
		go func() {
			wg.Add(1)
			ticker := time.NewTicker(time.Second)
			for {
				select {
				case <-ctx.Done():
					wg.Done()
					return
				case <-ticker.C:
					atomic.AddInt32(&count, 1)
				}
			}
		}()
	}

	wg.Wait()
	t.Logf("total count:%d \n", count)
}

func TestGetUserIntervalAndTimeIntervalV2(t *testing.T) {
	count := int32(0)
	userNum := 10000
	msgNum := 10000
	tm := time.Second

	userInv, timeInv := GetUserIntervalAndTimeInterval(userNum, msgNum, tm)
	t.Logf("userInterval %v, timeInterval: %v", userInv.String(), timeInv.String())

	ctx, closer := context.WithTimeout(context.Background(), time.Second*10)
	defer closer()

	wg := sync.WaitGroup{}
	for i := 0; i < userNum; i++ {
		go func(times int) {
			wg.Add(1)
			if times > 0 {
				inv := time.Duration(int64(userInv) * int64(times))
				fmt.Printf("userInv: %s\n", inv.String())
				timer := time.NewTimer(inv)
				select {
				case <-timer.C:
				case <-ctx.Done():
					wg.Done()
					return
				}
			}
			ticker := time.NewTicker(timeInv)
			for {
				select {
				case <-ctx.Done():
					wg.Done()
					return
				case <-ticker.C:
					atomic.AddInt32(&count, 1)
				}
			}
		}(i)
	}

	wg.Wait()
	t.Logf("total count:%d \n", count)
}
