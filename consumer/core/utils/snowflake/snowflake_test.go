package snowflake

import "testing"

func TestReverse(t *testing.T) {
	SnowflakeCore, _ := NewSnowflake(1)
	traceID := SnowflakeCore.TraceID()
	t.Log(traceID)
}
