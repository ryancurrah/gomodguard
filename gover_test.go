package gomodguard

import "testing"

func TestIsOkGoVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		constraint string
		goVersion  string
		ok         bool
		wantErr    bool
	}{
		{
			name:       "ok",
			constraint: "> 1.18",
			goVersion:  "1.19",
			ok:         true,
			wantErr:    false,
		},
		{
			name:       "not ok",
			constraint: "< 1.18",
			goVersion:  "1.19",
			ok:         false,
			wantErr:    false,
		},
		{
			name:       "error",
			constraint: "< 1.18",
			goVersion:  "test",
			ok:         false,
			wantErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ok, err := IsOkGoVersion(test.constraint, test.goVersion)
			if ok != test.ok {
				t.Fatalf("got %v want %v", ok, test.ok)
			}
			if (err != nil) != test.wantErr {
				t.Fatalf("got %v want %v", ok, test.ok)
			}
		})
	}
}
