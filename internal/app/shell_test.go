package app

import "testing"

func TestIsInteractiveCommand(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{name: "blocks vim", command: "vim", want: true},
		{name: "blocks sudo vim", command: "sudo vim test.txt", want: true},
		{name: "blocks less in pipeline", command: "cat README.md | less", want: true},
		{name: "allows ls", command: "ls -la", want: false},
		{name: "allows find pipe grep", command: "find . -name '*.go' | grep shell", want: false},
		{name: "allows env assignment with ls", command: "LC_ALL=C ls", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := isInteractiveCommand(tc.command)
			if got != tc.want {
				t.Fatalf("isInteractiveCommand(%q) = %v, want %v", tc.command, got, tc.want)
			}
		})
	}
}
