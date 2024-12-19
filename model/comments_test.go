package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	cases := []struct {
		name    string
		payload *CreateCommentPayload
		fail    bool
	}{
		{
			name: "正常系",
			payload: &CreateCommentPayload{
				ItemID:  1,
				UserID:  "user1",
				Comment: "comment1",
			},
			fail: false,
		},
		{
			name: "異常系: ItemIDが存在しない",
			payload: &CreateCommentPayload{
				UserID:  "user1",
				Comment: "comment1",
			},
			fail: true,
		},
		{
			name: "異常系: UserIDが存在しない",
			payload: &CreateCommentPayload{
				ItemID:  1,
				Comment: "comment1",
			},
			fail: true,
		},
		{
			name: "異常系: Commentが存在しない",
			payload: &CreateCommentPayload{
				ItemID: 1,
				UserID: "user1",
			},
			fail: true,
		},
	}

	assert := assert.New(t)
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateComment(tt.payload)
			if tt.fail {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
