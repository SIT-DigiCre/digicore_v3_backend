package utils

import "testing"

func TestCalculateSchoolGrade(t *testing.T) {
	tests := []struct {
		name          string
		studentNumber string
		schoolYear    int
		want          int
	}{
		{
			name:          "通常ケース",
			studentNumber: "aa23001",
			schoolYear:    2025,
			want:          3,
		},
		{
			name:          "院進ケース",
			studentNumber: "ma23001",
			schoolYear:    2025,
			want:          7,
		},
		{
			name:          "不正フォーマットは1年生フォールバック",
			studentNumber: "invalid",
			schoolYear:    2025,
			want:          1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateSchoolGrade(tt.studentNumber, tt.schoolYear)
			if got != tt.want {
				t.Fatalf("CalculateSchoolGrade() = %d, want %d", got, tt.want)
			}
		})
	}
}
