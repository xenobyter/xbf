package fs

import (
	"reflect"
	"testing"
)

func TestSortFiles(t *testing.T) {
	input := []FileInfo{
		{Name: "zebra.txt", IsDir: false},
		{Name: "Beta_dir", IsDir: true},
		{Name: "apple.txt", IsDir: false},
		{Name: "alpha_dir", IsDir: true},
		{Name: "B.txt", IsDir: false},
		{Name: "a.txt", IsDir: false},
	}

	expected := []FileInfo{
		{Name: "alpha_dir", IsDir: true},
		{Name: "Beta_dir", IsDir: true},
		{Name: "a.txt", IsDir: false},
		{Name: "apple.txt", IsDir: false},
		{Name: "B.txt", IsDir: false},
		{Name: "zebra.txt", IsDir: false},
	}

	sortFiles(input)

	for i := range expected {
		if !reflect.DeepEqual(input[i], expected[i]) {
			t.Errorf("at index %d: expected %+v, got %+v", i, expected[i], input[i])
		}
	}
}

func TestFilterFiles(t *testing.T) {
	files := []FileInfo{
		{Name: "Documents"},
		{Name: "report.txt"},
		{Name: "REPORT.md"},
	}

	got := filterFiles(files, "port")
	if len(got) != 2 || got[0].Name != "report.txt" || got[1].Name != "REPORT.md" {
		t.Fatalf("filterFiles() = %#v, want both case-insensitive matches", got)
	}

	if got := filterFiles(files, "missing"); len(got) != 0 {
		t.Fatalf("filterFiles() returned %d entries for a non-matching query", len(got))
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		size int64
		want string
	}{
		{size: 0, want: "0 B"},
		{size: 1023, want: "1023 B"},
		{size: 1024, want: "1.0 KB"},
		{size: 1024 * 1024, want: "1.0 MB"},
	}

	for _, test := range tests {
		if got := FormatFileSize(test.size); got != test.want {
			t.Errorf("FormatFileSize(%d) = %q, want %q", test.size, got, test.want)
		}
	}
}

func TestFilterHidden(t *testing.T) {
	files := []FileInfo{
		{Name: ".git", IsDir: true},
		{Name: ".gitignore", IsDir: false},
		{Name: "main.go", IsDir: false},
		{Name: ".env", IsDir: false},
		{Name: "src", IsDir: true},
	}

	expected := []FileInfo{
		{Name: "main.go", IsDir: false},
		{Name: "src", IsDir: true},
	}

	result := filterHidden(files)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("filterHidden() = %+v, expected %+v", result, expected)
	}

	// Also test file with IsHidden: true without leading dot (e.g. Windows hidden attribute)
	windowsHidden := []FileInfo{
		{Name: "desktop.ini", IsHidden: true},
		{Name: "normal.txt", IsHidden: false},
	}
	resWin := filterHidden(windowsHidden)
	if len(resWin) != 1 || resWin[0].Name != "normal.txt" {
		t.Errorf("expected desktop.ini with IsHidden=true to be filtered, got: %+v", resWin)
	}

	// Empty list
	if len(filterHidden(nil)) != 0 {
		t.Errorf("filterHidden(nil) should be empty")
	}

	// No hidden files
	noHidden := []FileInfo{{Name: "a.txt"}, {Name: "b.txt"}}
	if len(filterHidden(noHidden)) != 2 {
		t.Errorf("filterHidden(noHidden) should keep all files")
	}
}
