package utils

import "testing"

func TestURLFromHandle(t *testing.T) {
	type args struct {
		input   string
		siteURL string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "input is already a URL https",
			args: args{
				input:   "https://foo.com",
				siteURL: "https://bar.com",
			},
			want: "https://foo.com",
		},
		{
			name: "input is already a URL http",
			args: args{
				input:   "http://foo.com",
				siteURL: "https://bar.com",
			},
			want: "http://foo.com",
		},
		{
			name: "input is not a URL",
			args: args{
				input:   "foo",
				siteURL: "https://foo.com",
			},
			want: "https://foo.com/foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := URLFromHandle(tt.args.input, tt.args.siteURL); got != tt.want {
				t.Errorf("URLFromHandle() = %v, want %v", got, tt.want)
			}
		})
	}
}
