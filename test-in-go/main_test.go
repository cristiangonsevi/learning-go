package main

import "testing"

func TestSum(t *testing.T) {
	testCases := []struct {
		a    int
		b    int
		want int
	}{
		{
			a:    1,
			b:    2,
			want: 3,
		},
		{
			a:    3,
			b:    4,
			want: 8,
		},
	}
	for _, tt := range testCases {
		t.Run("testing", func(t *testing.T) {
			result := sum(tt.a, tt.b)

			if tt.want != result {
				t.Errorf("Error expeted %v but recieved %v", tt.want, result)
			}
		})
	}
}
