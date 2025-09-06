package commands

import (
	"testing"
)

func TestValidateEnodeURL(t *testing.T) {
	// Test cases for valid enode URLs
	validCases := []struct {
		name string
		url  string
	}{
		{
			name: "Valid enode URL",
			url:  "enode://9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d@192.168.1.100:30303",
		},
		{
			name: "Valid enode URL with IPv6",
			url:  "enode://9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d@[2001:db8::1]:30303",
		},
	}

	for _, tc := range validCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateEnodeURL(tc.url)
			if err != nil {
				t.Errorf("Expected no error for valid URL, got: %v", err)
			}
		})
	}

	// Test cases for invalid enode URLs
	invalidCases := []struct {
		name        string
		url         string
		expectedErr string
	}{
		{
			name:        "Missing enode prefix",
			url:         "9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d@192.168.1.100:30303",
			expectedErr: "enode URL must start with 'enode://'",
		},
		{
			name:        "Invalid format - missing @",
			url:         "enode://9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d",
			expectedErr: "invalid enode URL format: expected 'enode://<public-key>@<ip>:<port>'",
		},
		{
			name:        "Invalid public key - odd length",
			url:         "enode://9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6@192.168.1.100:30303",
			expectedErr: "공개 키는 128자 길이의 16진수 문자열이어야 합니다 (현재 길이: 127)",
		},
		{
			name:        "Invalid public key - non-hex character",
			url:         "enode://9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6g@192.168.1.100:30303",
			expectedErr: "공개 키는 16진수 문자열이어야 합니다",
		},
		{
			name:        "Invalid IP address",
			url:         "enode://9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d@999.999.999.999:30303",
			expectedErr: "IP 주소 형식이 올바르지 않습니다: 999.999.999.999",
		},
		{
			name:        "Invalid port - not a number",
			url:         "enode://9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d@192.168.1.100:abc",
			expectedErr: "포트 번호는 1-65535 범위의 숫자여야 합니다: abc",
		},
		{
			name:        "Invalid port - out of range",
			url:         "enode://9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d@192.168.1.100:99999",
			expectedErr: "포트 번호는 1-65535 범위의 숫자여야 합니다: 99999",
		},
	}

	for _, tc := range invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateEnodeURL(tc.url)
			if err == nil {
				t.Errorf("Expected error for invalid URL, but got none")
			} else if err.Error() != tc.expectedErr {
				t.Errorf("Expected error message '%s', but got '%s'", tc.expectedErr, err.Error())
			}
		})
	}
}

func TestValidatePublicKey(t *testing.T) {
	// Test valid public key
	validKey := "9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d"
	err := validatePublicKey(validKey)
	if err != nil {
		t.Errorf("Expected no error for valid public key, got: %v", err)
	}

	// Test invalid public key - odd length
	oddLengthKey := "9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6"
	err = validatePublicKey(oddLengthKey)
	if err == nil {
		t.Errorf("Expected error for odd length public key, but got none")
	}

	// Test invalid public key - non-hex character
	nonHexKey := "9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6d5e4f3e2d1c0b9a8b7c6g"
	err = validatePublicKey(nonHexKey)
	if err == nil {
		t.Errorf("Expected error for non-hex public key, but got none")
	}
}

func TestValidateIP(t *testing.T) {
	// Test valid IP addresses
	validIPs := []string{"192.168.1.1", "10.0.0.1", "2001:db8::1"}
	for _, ip := range validIPs {
		err := validateIP(ip)
		if err != nil {
			t.Errorf("Expected no error for valid IP %s, got: %v", ip, err)
		}
	}

	// Test invalid IP address
	invalidIP := "999.999.999.999"
	err := validateIP(invalidIP)
	if err == nil {
		t.Errorf("Expected error for invalid IP %s, but got none", invalidIP)
	}
}

func TestValidatePort(t *testing.T) {
	// Test valid ports
	validPorts := []string{"30303", "8545", "1", "65535"}
	for _, port := range validPorts {
		err := validatePort(port)
		if err != nil {
			t.Errorf("Expected no error for valid port %s, got: %v", port, err)
		}
	}

	// Test invalid ports
	invalidPorts := []struct {
		port string
		err  string
	}{
		{"abc", "포트 번호는 1-65535 범위의 숫자여야 합니다: abc"},
		{"99999", "포트 번호는 1-65535 범위의 숫자여야 합니다: 99999"},
		{"0", "포트 번호는 1-65535 범위의 숫자여야 합니다: 0"},
	}

	for _, tc := range invalidPorts {
		err := validatePort(tc.port)
		if err == nil {
			t.Errorf("Expected error for invalid port %s, but got none", tc.port)
		} else if err.Error() != tc.err {
			t.Errorf("Expected error message '%s', but got '%s'", tc.err, err.Error())
		}
	}
}