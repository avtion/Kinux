package services

import (
	"Kinux/core/web/models"
	"context"
	"testing"
	"time"
)

func TestNewMissionController(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{ctx: context.Background()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := NewMissionController(tt.args.ctx)

			ms, err := models.ListMissions(tt.args.ctx, "", nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(ms) == 0 {
				t.Fatal("no mission")
			}

			acs, err := models.ListAccounts(tt.args.ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(acs) == 0 {
				t.Fatal("account is nil")
			}

			if err = mc.SetAc(acs[0]).SetMission(ms[0]).NewDeployment(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func Test_MissionController_ResetAndWatch(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{ctx: context.Background()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(tt.args.ctx, 5*time.Minute)
			defer cancel()
			mc := NewMissionController(ctx)

			ms, err := models.ListMissions(tt.args.ctx, "", nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(ms) == 0 {
				t.Fatal("no mission")
			}

			acs, err := models.ListAccounts(tt.args.ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(acs) == 0 {
				t.Fatal("account is nil")
			}

			mc = mc.SetAc(acs[0]).SetMission(ms[0])

			errCh := mc.WatchDeploymentToReady("")

			if err = mc.ResetMission(""); err != nil {
				t.Fatal(err)
			}
			select {
			case err := <-errCh:
				if err != nil {
					t.Fatal("dp重置失败", err)
				}
				t.Log("dp重置成功")
			}
		})
	}
}
