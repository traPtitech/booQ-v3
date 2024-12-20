package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	PrepareTestDatabase()

	cases := []struct {
		name    string
		payload *CreateCommentPayload
		ok      bool
	}{
		{
			name: "正常系",
			payload: &CreateCommentPayload{
				ItemID:  1,
				UserID:  "user1",
				Comment: "comment1",
			},
			ok: true,
		},
		{
			name: "異常系: ItemIDが存在しない",
			payload: &CreateCommentPayload{
				UserID:  "user1",
				Comment: "comment1",
			},
			ok: false,
		},
		{
			name: "異常系: UserIDが存在しない",
			payload: &CreateCommentPayload{
				ItemID:  1,
				Comment: "comment1",
			},
			ok: false,
		},
		{
			name: "異常系: Commentが存在しない",
			payload: &CreateCommentPayload{
				ItemID: 1,
				UserID: "user1",
			},
			ok: false,
		},
	}

	assert := assert.New(t)
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateComment(tt.payload)
			if tt.ok {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
