package neko

import "testing"

func TestParseRouteTokenDoesntSpan(t *testing.T) {
	_, err := ParseRoute("/v1/api/d{id}", false)
	if err == nil {
		t.Errorf("Parser should have detected bad token")
	}
}

func TestParseRouteTokenMalformed(t *testing.T) {
	_, err := ParseRoute("/v1/api/{id-one}", false)
	if err == nil {
		t.Errorf("Parser should have detected malformed token")
	}
}

func TestParseRouteGoodToken(t *testing.T) {
	_, err := ParseRoute("/v1/api/{id}", false)
	if err != nil {
		t.Error(err)
	}
}

func TestParseRouteGoodTokenWithType(t *testing.T) {
	_, err := ParseRoute("/v1/api/{id:int}", false)
	if err != nil {
		t.Error(err)
	}
}
