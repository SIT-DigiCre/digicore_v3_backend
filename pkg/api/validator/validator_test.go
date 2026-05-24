package validator

import (
	"testing"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
)

func TestPhoneNumberValidation(t *testing.T) {
	tests := []struct {
		name                  string
		phoneNumber           string
		parentCellphoneNumber string
		parentHomephoneNumber *string
		expectErr             bool
	}{
		{
			name:                  "既存形式（日本形式）で通過",
			phoneNumber:           "09012345678",
			parentCellphoneNumber: "09087654321",
			parentHomephoneNumber: nil,
			expectErr:             false,
		},
		{
			name:                  "新形式（E.164）で通過",
			phoneNumber:           "+819011111111",
			parentCellphoneNumber: "+819087654321",
			parentHomephoneNumber: nil,
			expectErr:             false,
		},
		{
			name:                  "混合形式で通過",
			phoneNumber:           "09012345678",
			parentCellphoneNumber: "+819087654320",
			parentHomephoneNumber: nil,
			expectErr:             false,
		},
		{
			name:                  "固定電話なしで通過",
			phoneNumber:           "09012345678",
			parentCellphoneNumber: "09087654321",
			parentHomephoneNumber: nil,
			expectErr:             false,
		},
		{
			name:                  "固定電話ありで通過",
			phoneNumber:           "09012345678",
			parentCellphoneNumber: "09087654321",
			parentHomephoneNumber: stringPtr("0312345678"),
			expectErr:             false,
		},
		{
			name:                  "固定電話E.164形式で通過",
			phoneNumber:           "09012345678",
			parentCellphoneNumber: "09087654321",
			parentHomephoneNumber: stringPtr("+81312345678"),
			expectErr:             false,
		},
		{
			name:                  "米国番号で通過",
			phoneNumber:           "+18002752273",
			parentCellphoneNumber: "+14155552671",
			parentHomephoneNumber: nil,
			expectErr:             false,
		},
		{
			name:                  "米国番号と日本番号の混合で通過",
			phoneNumber:           "09012345678",
			parentCellphoneNumber: "+14155552671",
			parentHomephoneNumber: nil,
			expectErr:             false,
		},
		{
			name:                  "携帯が9桁で失敗",
			phoneNumber:           "090123456",
			parentCellphoneNumber: "09087654321",
			parentHomephoneNumber: nil,
			expectErr:             true,
		},
		{
			name:                  "携帯が12桁で失敗",
			phoneNumber:           "090123456789",
			parentCellphoneNumber: "09087654321",
			parentHomephoneNumber: nil,
			expectErr:             true,
		},
		{
			name:                  "無効な国コード付き番号で失敗",
			phoneNumber:           "+999999999999",
			parentCellphoneNumber: "09087654321",
			parentHomephoneNumber: nil,
			expectErr:             true,
		},
		{
			name:                  "空文字列で失敗",
			phoneNumber:           "",
			parentCellphoneNumber: "09087654321",
			parentHomephoneNumber: nil,
			expectErr:             true,
		},
		{
			name:                  "親携帯が空で失敗",
			phoneNumber:           "09012345678",
			parentCellphoneNumber: "",
			parentHomephoneNumber: nil,
			expectErr:             true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := api.ReqPutUserMePrivate{
				Address:               "住所",
				FirstName:             "太郎",
				FirstNameKana:         "タロウ",
				IsMale:                true,
				LastName:              "山田",
				LastNameKana:          "ヤマダ",
				ParentAddress:         "親の住所",
				ParentCellphoneNumber: tc.parentCellphoneNumber,
				ParentFirstName:       stringPtr("親太郎"),
				ParentHomephoneNumber: tc.parentHomephoneNumber,
				ParentLastName:        stringPtr("親山田"),
				PhoneNumber:           tc.phoneNumber,
			}

			err := Validate(req)
			if tc.expectErr {
				if err == nil {
					t.Fatal("エラーが発生することを期待していましたが、エラーはありませんでした")
				}
			} else {
				if err != nil {
					t.Fatalf("エラーが発生しませんことを期待していましたが、エラーが発生しました: %v", err)
				}
			}
		})
	}
}

// ヘルパー関数
func stringPtr(s string) *string {
	return &s
}
