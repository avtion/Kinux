package models

import (
	"context"
	"testing"
)

func TestCrateOrUpdateMission(t *testing.T) {
	dps, err := ListDeployment(context.Background(), "centos", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(dps) == 0 {
		t.Fatal("no dp")
	}
	dp := dps[0]
	type args struct {
		ctx  context.Context
		name string
		dp   *Deployment
		opts []MissionBuildOpt
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				ctx:  context.Background(),
				name: "centos",
				dp:   dp,
				opts: []MissionBuildOpt{
					MissionOptDesc("测试描述"),
					//MissionOptVnc("centos", "6777"),
					MissionOptDeployment("bash", "centos", []string{"centos"}),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotM, err := CrateOrUpdateMission(tt.args.ctx, tt.args.name, tt.args.dp, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CrateOrUpdateMission() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(gotM)
		})
	}
}

func TestEditMissionCheckpoints(t *testing.T) {
	type args struct {
		ctx         context.Context
		missionID   uint
		checkpoints []struct {
			CheckpointID    uint
			Percent         uint
			Priority        int
			TargetContainer string
		}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				ctx:       context.Background(),
				missionID: 1,
				checkpoints: []struct {
					CheckpointID    uint
					Percent         uint
					Priority        int
					TargetContainer string
				}{
					{
						CheckpointID:    1,
						Percent:         0,
						Priority:        0,
						TargetContainer: "",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EditMissionCheckpoints(tt.args.ctx, tt.args.missionID, tt.args.checkpoints...); (err != nil) != tt.wantErr {
				t.Errorf("EditMissionCheckpoints() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
