package notify

import (
	"testing"
)

func TestNotificator(t *testing.T) {
	// これ.envからtokenが読み込めないっぽくて
	if err := Notificator(); err != nil {
		t.Errorf("discordgo: errorを返すべきではない: %v", err)
	}
}
