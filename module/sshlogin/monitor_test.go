package sshlogin

import (
	"testing"
)

func TestParseSSHLogLine(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		wantUser   string
		wantHost   string
		wantPort   int
		wantStat   string
		wantAuth   string
		wantInval  int
		wantReason string
	}{
		{
			name:       "accepted password",
			line:       "Nov 15 22:00:01 myhost sshd[12345]: Accepted password for root from 192.168.1.100 port 54321 ssh2",
			wantUser:   "root",
			wantHost:   "192.168.1.100",
			wantPort:   54321,
			wantStat:   "success",
			wantAuth:   "password",
			wantInval:  0,
			wantReason: "",
		},
		{
			name:       "accepted publickey",
			line:       "Nov 15 22:00:02 myhost sshd[12346]: Accepted publickey for ubuntu from 10.0.0.1 port 52123 ssh2: RSA SHA256:abc123xyz",
			wantUser:   "ubuntu",
			wantHost:   "10.0.0.1",
			wantPort:   52123,
			wantStat:   "success",
			wantAuth:   "publickey",
			wantInval:  0,
			wantReason: "",
		},
		{
			name:       "failed password existing user",
			line:       "Nov 15 22:00:03 myhost sshd[12347]: Failed password for root from 192.168.1.200 port 43210 ssh2",
			wantUser:   "root",
			wantHost:   "192.168.1.200",
			wantPort:   43210,
			wantStat:   "failed",
			wantAuth:   "password",
			wantInval:  0,
			wantReason: "bad_password",
		},
		{
			name:       "failed password invalid user",
			line:       "Nov 15 22:00:04 myhost sshd[12348]: Failed password for invalid user admin from 192.168.1.201 port 43211 ssh2",
			wantUser:   "admin",
			wantHost:   "192.168.1.201",
			wantPort:   43211,
			wantStat:   "failed",
			wantAuth:   "password",
			wantInval:  1,
			wantReason: "invalid_user",
		},
		{
			name:       "failed publickey",
			line:       "Nov 15 22:00:05 myhost sshd[12349]: Failed publickey for dev from 172.16.0.2 port 60000 ssh2",
			wantUser:   "dev",
			wantHost:   "172.16.0.2",
			wantPort:   60000,
			wantStat:   "failed",
			wantAuth:   "publickey",
			wantInval:  0,
			wantReason: "bad_publickey",
		},
		{
			name:       "pam authentication failure fallback",
			line:       "Nov 15 22:00:06 myhost sshd[12350]: pam_unix(sshd:auth): authentication failure; logname= uid=0 euid=0 tty=ssh ruser= rhost=1.2.3.4  user=tester",
			wantUser:   "tester",
			wantHost:   "1.2.3.4",
			wantPort:   0,
			wantStat:   "failed",
			wantAuth:   "pam",
			wantInval:  0,
			wantReason: "pam_failure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := ParseSSHLogLine(tt.line)
			if ev == nil {
				t.Fatalf("expected non-nil event for line: %s", tt.line)
			}
			if ev.Username != tt.wantUser {
				t.Errorf("Username = %s, want %s", ev.Username, tt.wantUser)
			}
			if ev.RHost != tt.wantHost {
				t.Errorf("RHost = %s, want %s", ev.RHost, tt.wantHost)
			}
			if ev.Port != tt.wantPort {
				t.Errorf("Port = %d, want %d", ev.Port, tt.wantPort)
			}
			if ev.Status != tt.wantStat {
				t.Errorf("Status = %s, want %s", ev.Status, tt.wantStat)
			}
			if ev.AuthType != tt.wantAuth {
				t.Errorf("AuthType = %s, want %s", ev.AuthType, tt.wantAuth)
			}
			if ev.IsInvalid != tt.wantInval {
				t.Errorf("IsInvalid = %d, want %d", ev.IsInvalid, tt.wantInval)
			}
			if ev.FailReason != tt.wantReason {
				t.Errorf("FailReason = %s, want %s", ev.FailReason, tt.wantReason)
			}
		})
	}
}
