package diff

import (
	"testing"

	r3 "github.com/r3labs/diff"
)

type SampleObj struct {
	StringValue string
	IntValue    int
	StringArray []string
	StringMap   map[string]string
}

func TestGetDiffChangelog(t *testing.T) {
	type args struct {
		oldObj interface{}
		newObj interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *r3.Changelog
		wantErr bool
	}{
		{
			name: "success get diff change log",
			args: args{
				oldObj: SampleObj{StringValue: "string1", IntValue: 1, StringArray: []string{"strarray1"}, StringMap: map[string]string{"strmap": "strmapvalue1"}},
				newObj: SampleObj{StringValue: "string2", IntValue: 2, StringArray: []string{"strarray2"}, StringMap: map[string]string{"strmap": "strmapvalue2"}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetDiffChangelog(tt.args.oldObj, tt.args.newObj)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDiffChangelog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetDiffString(t *testing.T) {
	type args struct {
		oldObj interface{}
		newObj interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "success get diff string",
			args: args{
				oldObj: SampleObj{StringValue: "string1", IntValue: 1, StringArray: []string{"strarray1"}, StringMap: map[string]string{"strmap": "strmapvalue1"}},
				newObj: SampleObj{StringValue: "string2", IntValue: 2, StringArray: []string{"strarray2"}, StringMap: map[string]string{"strmap": "strmapvalue2"}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetDiffString(tt.args.oldObj, tt.args.newObj)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDiffString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
