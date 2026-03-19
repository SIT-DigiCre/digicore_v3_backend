package admin

import (
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

type fakeTransactionClient struct {
	selectFunc func(dest interface{}, queryPath string, params interface{}) error
	execFunc   func(queryPath string, params interface{}, generateId bool) (sql.Result, error)
}

func (f *fakeTransactionClient) Select(dest interface{}, queryPath string, params interface{}) error {
	if f.selectFunc != nil {
		return f.selectFunc(dest, queryPath, params)
	}
	return nil
}

func (f *fakeTransactionClient) Exec(queryPath string, params interface{}, generateId bool) (sql.Result, error) {
	if f.execFunc != nil {
		return f.execFunc(queryPath, params, generateId)
	}
	return nil, nil
}

func (f *fakeTransactionClient) GetId() (string, error) {
	return "", nil
}

func (f *fakeTransactionClient) DuplicateUpdate(string, string, interface{}) (sql.Result, bool, error) {
	return nil, false, nil
}

func TestPutAdminSchoolGradeUpdatesEachTarget(t *testing.T) {
	client := &fakeTransactionClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			if queryPath != "sql/admin/select_school_grade_update_targets.sql" {
				t.Fatalf("unexpected query path: %s", queryPath)
			}

			rows := dest.(*[]struct {
				UserId             string `db:"user_id"`
				StudentNumber      string `db:"student_number"`
				ApprovedGradeDiffs int    `db:"approved_grade_diffs"`
			})
			*rows = []struct {
				UserId             string `db:"user_id"`
				StudentNumber      string `db:"student_number"`
				ApprovedGradeDiffs int    `db:"approved_grade_diffs"`
			}{
				{UserId: "user-1", StudentNumber: "aa25001", ApprovedGradeDiffs: 0},
				{UserId: "user-2", StudentNumber: "m250001", ApprovedGradeDiffs: -1},
			}
			return nil
		},
	}

	type updateParams struct {
		UserId      string
		SchoolGrade int
	}
	var updates []updateParams
	client.execFunc = func(queryPath string, params interface{}, generateId bool) (sql.Result, error) {
		if queryPath != "sql/admin/update_user_profile_school_grade.sql" {
			t.Fatalf("unexpected query path: %s", queryPath)
		}
		value := reflect.ValueOf(params).Elem()
		var rawUpdates []struct {
			UserId      string `json:"userId"`
			SchoolGrade int    `json:"schoolGrade"`
		}
		if err := json.Unmarshal([]byte(value.FieldByName("UpdatesJSON").String()), &rawUpdates); err != nil {
			t.Fatalf("failed to unmarshal updates json: %v", err)
		}
		for _, update := range rawUpdates {
			updates = append(updates, updateParams{
				UserId:      update.UserId,
				SchoolGrade: update.SchoolGrade,
			})
		}
		return nil, nil
	}

	res, err := PutAdminSchoolGrade(echo.New().NewContext(nil, nil), client)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Success {
		t.Fatal("expected success response")
	}
	if len(updates) != 2 {
		t.Fatalf("expected 2 updates, got %d", len(updates))
	}
	expectedFirstGrade, calcErr := utils.CalculateSchoolGradeFromStudentNumber("aa25001")
	if calcErr != nil {
		t.Fatalf("unexpected error: %v", calcErr)
	}
	if updates[0].UserId != "user-1" || updates[0].SchoolGrade != expectedFirstGrade {
		t.Fatalf("unexpected first update: %+v", updates[0])
	}
	expectedGraduateGrade, calcErr := utils.CalculateSchoolGradeFromStudentNumber("m250001")
	if calcErr != nil {
		t.Fatalf("unexpected error: %v", calcErr)
	}
	expectedGraduateGrade -= 1
	if updates[1].UserId != "user-2" || updates[1].SchoolGrade != expectedGraduateGrade {
		t.Fatalf("unexpected second update: %+v", updates[1])
	}
}

func TestPutAdminSchoolGradeReturnsErrorWhenSelectFails(t *testing.T) {
	client := &fakeTransactionClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			return errors.New("select failed")
		},
	}

	_, err := PutAdminSchoolGrade(echo.New().NewContext(nil, nil), client)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Message != "学年の更新に失敗しました" {
		t.Fatalf("unexpected message: %s", err.Message)
	}
}

func TestPutAdminSchoolGradeReturnsErrorWhenCalculationFails(t *testing.T) {
	client := &fakeTransactionClient{
		selectFunc: func(dest interface{}, queryPath string, params interface{}) error {
			rows := dest.(*[]struct {
				UserId             string `db:"user_id"`
				StudentNumber      string `db:"student_number"`
				ApprovedGradeDiffs int    `db:"approved_grade_diffs"`
			})
			*rows = []struct {
				UserId             string `db:"user_id"`
				StudentNumber      string `db:"student_number"`
				ApprovedGradeDiffs int    `db:"approved_grade_diffs"`
			}{
				{UserId: "user-1", StudentNumber: "abxx001", ApprovedGradeDiffs: 0},
			}
			return nil
		},
	}

	_, err := PutAdminSchoolGrade(echo.New().NewContext(nil, nil), client)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Message != "学年の計算に失敗しました" {
		t.Fatalf("unexpected message: %s", err.Message)
	}
}
