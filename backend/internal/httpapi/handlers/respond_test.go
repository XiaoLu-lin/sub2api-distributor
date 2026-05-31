package handlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/lhl/sub2api-distributor/backend/internal/distributor"
)

func TestMessageForError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status int
		err    error
		want   string
	}{
		{
			name:   "keeps known business error message",
			status: http.StatusBadRequest,
			err:    distributor.ErrWithdrawalAmountTooLarge,
			want:   "提现金额超过当前可申请金额",
		},
		{
			name:   "hides unknown internal error details",
			status: http.StatusInternalServerError,
			err:    errors.New("query withdraw summary: pq timeout"),
			want:   internalServerMessage,
		},
		{
			name:   "keeps explicit client error message for bad request",
			status: http.StatusBadRequest,
			err:    errors.New("请求参数不正确"),
			want:   "请求参数不正确",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := messageForError(tt.status, tt.err)
			if got != tt.want {
				t.Fatalf("messageForError(%d, %v) = %q, want %q", tt.status, tt.err, got, tt.want)
			}
		})
	}
}
