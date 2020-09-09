package main

import (
	"reflect"
	"testing"
)

func TestFormat(t *testing.T) {
	type args struct {
		m rawMedia
	}
	tests := []struct {
		name    string
		args    args
		want    rawMedia
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := format(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("format() got = %v, want %v", got, tt.want)
			}
		})
	}
}
