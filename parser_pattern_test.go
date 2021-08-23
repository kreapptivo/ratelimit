package ratelimit

import (
	"testing"
	"time"

	"aahframework.org/test.v0/assert"
)

func Test_simple_pattern_SECOND(t *testing.T) {
	l, _ := parse("1r/s")
	assert.Equal(t, time.Second, l.Per)
	assert.Equal(t, uint32(1), l.Max)
}

func Test_simple_pattern_MINUTE(t *testing.T) {
	l, _ := parse("1r/m")
	assert.Equal(t, time.Minute, l.Per)
	assert.Equal(t, uint32(1), l.Max)
}

func Test_simple_pattern_HOUR(t *testing.T) {
	l, _ := parse("1r/h")
	assert.Equal(t, time.Hour, l.Per)
	assert.Equal(t, uint32(1), l.Max)
}

func Test_simple_pattern_DAY(t *testing.T) {
	l, _ := parse("1r/d")
	assert.Equal(t, time.Hour*24, l.Per)
	assert.Equal(t, uint32(1), l.Max)
}

func Test_pattern_with_spammer_and_blocker(t *testing.T) {
	l, _ := parse("1r/s,spam:2,block:15d")
	assert.Equal(t, uint32(2), l.MaxToSpam)
	assert.Equal(t, 15*(time.Hour*24), l.Block)

	l2, _ := parse("1r/m,spam:5,block:12d")
	assert.Equal(t, uint32(5), l2.MaxToSpam)
	assert.Equal(t, 12*(time.Hour*24), l2.Block)
}

func Test_it_panics_when_invalid_spam_value(t *testing.T) {
	_, err := parse("1r/s,spam,block:15d")
	expected := "Can't parse value: spam"
	if err == nil || err.Error() != expected {
		t.Errorf("Error actual = %v, and Expected = %v.", err.Error(), expected)
	}
}

func Test_it_panics_when_invalid_spam_params(t *testing.T) {
	_, err := parse("1r/s,spam:12:3,block:15d")
	expected := "Can't parse value: spam:12:3"
	if err == nil || err.Error() != expected {
		t.Errorf("Error actual = %v, and Expected = %v.", err.Error(), expected)
	}
}

func Test_it_panics_when_invalid_module(t *testing.T) {
	_, err := parse("1r/s,spam:12,block:15d,fake:3")
	expected := "Unsupported module [fake] must be spam or block"
	if err == nil || err.Error() != expected {
		t.Errorf("Error actual = %v, and Expected = %v.", err.Error(), expected)
	}
}

func Test_block_durations(t *testing.T) {
	l1, _ := parse("1r/s,spam:12,block:3d")
	l2, _ := parse("1r/s,spam:12,block:2m")
	l3, _ := parse("1r/s,spam:12,block:3s")
	l4, _ := parse("1r/s,spam:12,block:3h")
	assert.Equal(t, 3*(24*time.Hour), l1.Block)
	assert.Equal(t, 3*time.Hour, l4.Block)
	assert.Equal(t, 2*time.Minute, l2.Block)
	assert.Equal(t, 3*time.Second, l3.Block)
}

func Test_rate_durations(t *testing.T) {
	l1, _ := parse("1r/s,spam:12,block:1d")
	l2, _ := parse("1r/m,spam:12,block:1d")
	l3, _ := parse("1r/h,spam:12,block:1d")
	l4, _ := parse("1r/d,spam:12,block:1d")
	assert.Equal(t, time.Second, l1.Per)
	assert.Equal(t, time.Minute, l2.Per)
	assert.Equal(t, time.Hour, l3.Per)
	assert.Equal(t, time.Hour*24, l4.Per)
}

func Test_it_panics_when_invalid_block_duration(t *testing.T) {
	_, err := parse("1r/s,spam:12,block:15w")
	expected := "Unsupported time duration [w] must be (d) for day, (h) for hour, (m) for minute or (s) for second."
	if err == nil || err.Error() != expected {
		t.Errorf("Error actual = %v, and Expected = %v.", err.Error(), expected)
	}
}
