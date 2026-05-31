package handlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/lhl/sub2api-distributor/backend/internal/distributor"
)

func TestOperatorMutationStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		err    error
		want   int
	}{
		{
			name: "not found stays 404",
			err:  distributor.ErrWithdrawalNotFound,
			want: http.StatusNotFound,
		},
		{
			name: "invalid transition stays 400",
			err:  distributor.ErrInvalidWithdrawalTransition,
			want: http.StatusBadRequest,
		},
		{
			name: "unexpected error becomes 500",
			err:  errors.New("driver timeout"),
			want: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := operatorMutationStatus(tt.err)
			if got != tt.want {
				t.Fatalf("operatorMutationStatus(%v) = %d, want %d", tt.err, got, tt.want)
			}
		})
	}
}
