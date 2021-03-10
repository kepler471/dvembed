package main

import (
	"reflect"
	"testing"
)

func TestFormat(t *testing.T) {
	type args struct {
		m media
	}
	tests := []struct {
		name    string
		args    args
		want    media
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compress(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compress() got = %v, want %v", got, tt.want)
			}
		})
	}
}
