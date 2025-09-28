package service

import "testing"

func TestNormalizeAmount(t *testing.T) {
	tests := []struct {
		name string
		opID int
		in   float64
		out  float64
	}{
		{"cash purchase negative", 1, 50.0, -50.0},
		{"installment purchase negative", 2, -23.5, -23.5},
		{"withdrawal negative", 3, 18.7, -18.7},
		{"payment positive", 4, -60.0, 60.0},
		{"payment positive 2", 4, 60.0, 60.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeAmount(tt.opID, tt.in)
			if got != tt.out {
				t.Fatalf("NormalizeAmount(%d, %v) = %v; want %v", tt.opID, tt.in, got, tt.out)
			}
		})
	}
}



