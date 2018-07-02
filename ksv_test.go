package main

import (
	"strings"
	"testing"
)

var encoded_yaml_data = `
apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  username: YWRtaW4=
  password: MWYyZDFlMmU2N2Rm
`

var decoded_yaml_data = `
apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  username: admin
  password: 1f2d1e2e67df
`

func TestEncode(t *testing.T) {
	s, err := encodeToBase64(strings.NewReader(decoded_yaml_data))
	if err != nil {
		t.Error("Can't parse yaml doc")
	}
	if s.APIVersion != "v1" {
		t.Error("Invalid API version")
	}
	if s.Kind != "Secret" {
		t.Error("Invalid Kind")
	}
	if s.Data["username"] != "YWRtaW4=" {
		t.Error("Invalid username value")
	}
	if s.Data["password"] != "MWYyZDFlMmU2N2Rm" {
		t.Error("Invalid password value")
	}
	if len(s.StringData) != 0 {
		t.Error("StringData should be empty")
	}
}

func TestDecode(t *testing.T) {
	s, err := decodeFromBase64(strings.NewReader(encoded_yaml_data), false)
	if err != nil {
		t.Error("Can't parse yaml doc")
	}
	if s.APIVersion != "v1" {
		t.Error("Invalid API version")
	}
	if s.Kind != "Secret" {
		t.Error("Invalid Kind")
	}
	if s.Data["username"] != "admin" {
		t.Error("Invalid username value")
	}
	if s.Data["password"] != "1f2d1e2e67df" {
		t.Error("Invalid password value")
	}
	if len(s.StringData) != 0 {
		t.Error("StringData should be empty")
	}
}

func TestDecodeToStringData(t *testing.T) {
	s, err := decodeFromBase64(strings.NewReader(encoded_yaml_data), true)
	if err != nil {
		t.Error("Can't parse yaml doc")
	}
	if s.APIVersion != "v1" {
		t.Error("Invalid API version")
	}
	if s.Kind != "Secret" {
		t.Error("Invalid Kind")
	}
	if len(s.Data) != 0 {
		t.Error("Data should be empty")
	}
	if len(s.StringData) != 2 {
		t.Error("Not enough keys in stringData")
	}
	if s.StringData["username"] != "admin" {
		t.Error("Invalid username value")
	}
	if s.StringData["password"] != "1f2d1e2e67df" {
		t.Error("Invalid password value")
	}
}
