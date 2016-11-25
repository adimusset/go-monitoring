package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// The nextState method wraps all the alerting logic
// Having such a pure function permits to write very simple tests

func TestNextStateUp(t *testing.T) {
	now := time.Date(2016, time.November, 19, 12, 0, 0, 0, time.UTC)
	objects := []Object{
		Object{Date: now.Add(-150 * time.Second)},
		Object{Date: now.Add(-121 * time.Second)},
		Object{Date: now.Add(-119 * time.Second)},
		Object{Date: now.Add(-90 * time.Second)},
	}
	i, overAverage, alert := nextState(now, objects, false, 3, 120)
	assert.Nil(t, alert)
	assert.False(t, overAverage)
	assert.Equal(t, 2, i)

	i, overAverage, alert = nextState(now, objects, false, 3, 180)
	expectedAlert := &Alert{Date: now, Average: 3, Up: true}
	assert.Equal(t, expectedAlert, alert)
	assert.True(t, overAverage)
	assert.Equal(t, 0, i)
}

func TestNextStateDown(t *testing.T) {
	now := time.Date(2016, time.November, 19, 12, 0, 0, 0, time.UTC)
	objects := []Object{
		Object{Date: now.Add(-150 * time.Second)},
		Object{Date: now.Add(-121 * time.Second)},
		Object{Date: now.Add(-119 * time.Second)},
		Object{Date: now.Add(-90 * time.Second)},
	}
	i, overAverage, alert := nextState(now, objects, true, 3, 180)

	assert.Nil(t, alert)
	assert.True(t, overAverage)
	assert.Equal(t, 0, i)

	i, overAverage, alert = nextState(now, objects, true, 3, 120)
	expectedAlert := &Alert{Date: now, Average: 3, Up: false}
	assert.Equal(t, expectedAlert, alert)
	assert.False(t, overAverage)
	assert.Equal(t, 2, i)
}
