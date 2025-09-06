package commands

import (
	"testing"
)

func TestShowLocalEnodeText(t *testing.T) {
	// This is a basic test to ensure the function doesn't panic
	// In a real implementation, we would mock the network manager
	
	err := showLocalEnodeText()
	if err != nil {
		t.Errorf("showLocalEnodeText() error = %v", err)
	}
}

func TestShowLocalEnodeJSON(t *testing.T) {
	// This is a basic test to ensure the function doesn't panic
	// In a real implementation, we would mock the network manager
	
	err := showLocalEnodeJSON()
	if err != nil {
		t.Errorf("showLocalEnodeJSON() error = %v", err)
	}
}