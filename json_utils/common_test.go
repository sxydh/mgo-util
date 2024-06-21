package json_utils

import "testing"

func TestToJsonStr(t *testing.T) {
	type args struct {
		p interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Normal",
			args: args{
				p: &struct {
					Id string `json:"id"`
				}{
					Id: "1",
				}},
			want: "{\"id\":\"1\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJsonStr(tt.args.p); got != tt.want {
				t.Errorf("ToJsonStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
