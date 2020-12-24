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
				name: "centos-测试2",
				dp:   dp,
				opts: []MissionBuildOpt{
					MissionOptNs("test"),
					MissionOptDesc("测试描述"),
					MissionOptVnc("centos", "6777"),
					MissionOptDeployment("bash", "centos", []string{"1", "2", "3"}),
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
