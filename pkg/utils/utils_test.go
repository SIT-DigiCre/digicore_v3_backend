package utils

import "testing"

func TestCalculateSchoolGradeFromStudentNumber(t *testing.T) {
	tests := []struct {
		name              string
		studentNumber     string
		currentSchoolYear int
		expected          int
		expectErr         bool
	}{
		{
			name:              "学部生の学年を計算できる",
			studentNumber:     "aa25001",
			currentSchoolYear: 2025,
			expected:          1,
		},
		{
			name:              "大学院修士は 4 年加算する",
			studentNumber:     "m025001",
			currentSchoolYear: 2025,
			expected:          5,
		},
		{
			name:              "大学院博士は 6 年加算する",
			studentNumber:     "N025001",
			currentSchoolYear: 2025,
			expected:          7,
		},
		{
			name:              "古い入学年度でも負数エラーにならず計算できる",
			studentNumber:     "aa99001",
			currentSchoolYear: 2025,
			expected:          -73,
		},
		{
			name:              "学籍番号が短すぎるとエラー",
			studentNumber:     "a25",
			currentSchoolYear: 2025,
			expectErr:         true,
		},
		{
			name:              "入学年度が数値でないとエラー",
			studentNumber:     "aaxx001",
			currentSchoolYear: 2025,
			expectErr:         true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := calculateSchoolGradeFromStudentNumber(tc.studentNumber, tc.currentSchoolYear)
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
