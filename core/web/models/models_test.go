package models

import (
	"Kinux/tools/cfg"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"testing"
)

// initDatabase
func init() {
	cfg.InitConfig(func(v *viper.Viper) {
		v.AddConfigPath("../../../")
	})
	if err := InitDatabaseConn(context.Background(),
		cfg.DefaultConfig.Database.Name, cfg.DefaultConfig.Database.Dsn); err != nil {
		logrus.Fatal(err)
	}
}

func TestInitDatabaseConn(t *testing.T) {
	type args struct {
		ctx context.Context
		dsn string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{

		{
			name: "test",
			args: args{
				ctx: context.Background(),
				dsn: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg.InitConfig(func(v *viper.Viper) {
				v.AddConfigPath("../../../")
			})
			if err := InitDatabaseConn(tt.args.ctx,
				cfg.DefaultConfig.Database.Name, cfg.DefaultConfig.Database.Dsn); (err != nil) != tt.wantErr {
				t.Errorf("InitDatabaseConn() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(GetGlobalDB())
		})
	}
}
