package cache

import "testing"

func Test_getPrefix(t *testing.T) {
	type args struct {
		prefix      string
		groupPrefix GroupType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{"", ""}, ""},
		{"empty group", args{"prefix:", ""}, "prefix:"},
		{"empty prefix", args{"", GroupTypeFUP}, "fup_"},
		{"both", args{"prefix:", GroupTypeFUP}, "prefix:fup_"},
		{"both with underscore", args{"prefix_", GroupTypeAuth}, "prefix_auth_"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPrefix(tt.args.prefix, tt.args.groupPrefix); got != tt.want {
				t.Errorf("getPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
