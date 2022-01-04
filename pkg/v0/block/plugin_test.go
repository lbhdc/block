package block

import (
	"testing"
	"time"
)

func TestPlugin(t *testing.T) {
	tests := []struct {
		name       string
		entrypoint string
		wait       time.Duration
	}{
		{
			name:       "plugins start and stop",
			entrypoint: "while true; do sleep 1; done",
			wait:       10 * time.Millisecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin := NewPlugin(tt.entrypoint)
			if err := plugin.Start(); err != nil {
				t.Errorf("failed to start plugin, err = %v", err)
			}
			<-time.After(tt.wait)
			if err := plugin.Kill(); err != nil {
				t.Errorf("failed to kill plugin, err = %v", err)
			}
		})
	}
}
