package helpers

import "testing"

func TestGetHash(t *testing.T) {
	cases := []string{
		"hello",
		"привет, мир!",
		"WlhsS01HVllRV2xQYVVwTFZqRlJhVXhEU21oaVIyTnBUMmxLU1ZWNlNURk9hVW81TG1WNVNuQmFRMGsyU1dwRmVrMTZZMmxNUTBveFl6SldlV0p0Um5SYVUwazJTVzFLY0dWdE9YVmFVMGx6U1cxc2FHUkRTVFpOVkZVMVRrUkpkMDlVV1hkTlEzZHBZMjA1YzFwVFNUWkpibFo2V2xoSmFXWlJMbHAyYTFsWmJubE5PVEk1UmswMFRsYzVYMmhUYVhNM1gzZ3pYemx5ZVcxelJFRjRPWGwxVDJOak1Vaz0=",
	}
	for _, v := range cases {
		result, err := GetHash(v)
		if err != nil {
			t.Errorf("Фраза '%s' выдала ошибку %s", v, err.Error())
		}
		if !CheckHash(v, result) {
			t.Errorf("Несовпадение: '%s', '%s'", v, result)
		}
	}
}

func TestEncodeString(t *testing.T) {
	cases := map[string]string{
		"hello": "aGVsbG8=",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZCI6IjEzMzciLCJ1c2VybmFtZSI6ImJpem9uZSIsImlhdCI6MTU5NDIwOTYwMCwicm9sZSI6InVzZXIifQ.ZvkYYnyM929FM4NW9_hSis7_x3_9rymsDAx9yuOcc1I": "ZXlKMGVYQWlPaUpLVjFRaUxDSmhiR2NpT2lKSVV6STFOaUo5LmV5SnBaQ0k2SWpFek16Y2lMQ0oxYzJWeWJtRnRaU0k2SW1KcGVtOXVaU0lzSW1saGRDSTZNVFU1TkRJd09UWXdNQ3dpY205c1pTSTZJblZ6WlhJaWZRLlp2a1lZbnlNOTI5Rk00Tlc5X2hTaXM3X3gzXzlyeW1zREF4OXl1T2NjMUk=",
	}
	for k, v := range cases {
		r := EncodeString(k)
		if r != v {
			t.Errorf("%s: получено '%s', ожидалось '%s'", k, r, v)
		}
	}
}

func TestDecodeString(t *testing.T) {
	cases := map[string]string{
		"aGVsbG8=":                     "hello",
		"0L/RgNC40LLQtdGCLCDQvNC40YAh": "привет, мир!",
		"ZXlKMGVYQWlPaUpLVjFRaUxDSmhiR2NpT2lKSVV6STFOaUo5LmV5SnBaQ0k2SWpFek16Y2lMQ0oxYzJWeWJtRnRaU0k2SW1KcGVtOXVaU0lzSW1saGRDSTZNVFU1TkRJd09UWXdNQ3dpY205c1pTSTZJblZ6WlhJaWZRLlp2a1lZbnlNOTI5Rk00Tlc5X2hTaXM3X3gzXzlyeW1zREF4OXl1T2NjMUk=": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZCI6IjEzMzciLCJ1c2VybmFtZSI6ImJpem9uZSIsImlhdCI6MTU5NDIwOTYwMCwicm9sZSI6InVzZXIifQ.ZvkYYnyM929FM4NW9_hSis7_x3_9rymsDAx9yuOcc1I",
	}
	for k, v := range cases {
		r, err := DecodeString(k)
		if err != nil {
			t.Errorf("%s: %v", k, err)
		}
		if r != v {
			t.Errorf("%s: получено '%s', ожидалось '%s'", k, r, v)
		}
	}
}
