package utils

import "testing"

func TestCalculateSchoolGradeFromStudentNumber(t *testing.T) {
	currentSchoolYear := GetSchoolYear()
	tests := []struct {
		name          string
		studentNumber string
		expected      int
		expectErr     bool
	}{
		{
			name:          "学部生の学年を計算できる",
			studentNumber: "aa25001",
			expected:      currentSchoolYear - 2000 - 25 + 1,
		},
		{
			name:          "大学院修士は 4 年加算する",
			studentNumber: "m250001",
			expected:      currentSchoolYear - 2000 - 25 + 1 + 4,
		},
		{
			name:          "大学院博士は 6 年加算する",
			studentNumber: "N250001",
			expected:      currentSchoolYear - 2000 - 25 + 1 + 6,
		},
		{
			name:          "古い入学年度でも負数エラーにならず計算できる",
			studentNumber: "aa99001",
			expected:      currentSchoolYear - 2000 - 99 + 1,
		},
		{
			name:          "学籍番号が短すぎるとエラー",
			studentNumber: "a25",
			expectErr:     true,
		},
		{
			name:          "入学年度が数値でないとエラー",
			studentNumber: "aaxx001",
			expectErr:     true,
		},
		{
			name:          "大学院生の入学年度が数値でないとエラー",
			studentNumber: "mxx0001",
			expectErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := CalculateSchoolGradeFromStudentNumber(tc.studentNumber)
			if tc.expectErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if actual != tc.expected {
				t.Fatalf("expected %d, got %d", tc.expected, actual)
			}
		})
	}
}
