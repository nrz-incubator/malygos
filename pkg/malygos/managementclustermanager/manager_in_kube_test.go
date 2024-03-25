package managementclustermanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getIDFromSecretName(t *testing.T) {
	type args struct {
		secretName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "test",
			args: args{
				"pouet",
			},
			want: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := getIDFromSecretName(tt.args.secretName)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
