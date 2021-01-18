package main

import (
	"reflect"
	"testing"
)

func Test_jsonCheckValues(t *testing.T) {
	tests := []struct {
		name       string
		json       string
		jsonCheck  string
		wantErr    error
		wantResult bool
	}{
		{
			"status equals UP",
			"{\"status\": \"UP\"}",
			"status == UP",
			nil,
			true,
		},
		{
			"status equals 1",
			"{\"status\": 1}",
			"status == 1",
			nil,
			false,
		},
		{
			"status not equals UP",
			"{\"status\": \"DOWN\"}",
			"status == UP",
			nil,
			false,
		},
		{
			"status missing",
			"{\"current\": \"UP\"}",
			"status == OK",
			nil,
			false,
		},
		{
			"status not equals PARTIAL",
			"{\"status\": \"UP\"}",
			"status!=PARTIAL",
			nil,
			true,
		},
		{
			"status equals DOWN",
			"{\"status\": \"DOWN\"}",
			"status != DOWN",
			nil,
			false,
		},
		{
			"status not equals missing",
			"{\"current\": \"UP\"}",
			"status!=OK",
			nil,
			false,
		},
		{
			"status in UP",
			"{\"current\": \"UP\"}",
			"status in UP,PARTIAL",
			nil,
			false,
		},
		{
			"status in PARTIAL",
			"{\"current\": \"PARTIAL\"}",
			"status in UP,PARTIAL",
			nil,
			false,
		},
		{
			"status not in UP, PARTIAL",
			"{\"current\": \"DOWN\"}",
			"status in UP,PARTIAL",
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonCheck.Set(tt.jsonCheck)
			got1, got := jsonCheckValues([]byte(tt.json))
			if !reflect.DeepEqual(got, tt.wantErr) {
				t.Errorf("jsonCheckValues() got = %v, want %v", got, tt.wantErr)
			}
			if got1 != tt.wantResult {
				t.Errorf("jsonCheckValues() got1 = %v, want %v", got1, tt.wantResult)
			}
		})
	}
}
