package distributor

import (
	"context"
	"errors"
	"testing"
)

type fakeProfileStore struct {
	upsertCalls int
}

func (s *fakeProfileStore) upsertProfile(context.Context, Profile) error {
	s.upsertCalls++
	return nil
}

type fakeAffiliateEnsurer struct {
	calls int
	err   error
}

func (e *fakeAffiliateEnsurer) ensure(context.Context, int64) error {
	e.calls++
	return e.err
}

type fakeTxController struct {
	commits   int
	rollbacks int
}

func (c *fakeTxController) commit() error {
	c.commits++
	return nil
}

func (c *fakeTxController) rollback() {
	c.rollbacks++
}

func TestUpsertActiveProfileRequiresAffiliateIdentityBeforeSuccess(t *testing.T) {
	store := &fakeProfileStore{}
	ensurer := &fakeAffiliateEnsurer{err: errors.New("affiliate query failed")}
	tx := &fakeTxController{}

	err := runProfileUpsertTransaction(
		context.Background(),
		tx.commit,
		tx.rollback,
		profileStoreFunc(store.upsertProfile),
		ensurer.ensure,
		Profile{UserID: 42, Status: "active", DisplayName: "demo"},
	)
	if err == nil {
		t.Fatalf("runProfileUpsertTransaction() error = nil, want non-nil")
	}
	if store.upsertCalls != 1 {
		t.Fatalf("runProfileUpsertTransaction() upsertCalls = %d, want 1", store.upsertCalls)
	}
	if ensurer.calls != 1 {
		t.Fatalf("runProfileUpsertTransaction() ensurer.calls = %d, want 1", ensurer.calls)
	}
	if tx.commits != 0 {
		t.Fatalf("runProfileUpsertTransaction() commits = %d, want 0", tx.commits)
	}
	if tx.rollbacks != 1 {
		t.Fatalf("runProfileUpsertTransaction() rollbacks = %d, want 1", tx.rollbacks)
	}
}

func TestUpsertActiveProfileCommitsAfterAffiliateIdentitySuccess(t *testing.T) {
	store := &fakeProfileStore{}
	ensurer := &fakeAffiliateEnsurer{}
	tx := &fakeTxController{}

	err := runProfileUpsertTransaction(
		context.Background(),
		tx.commit,
		tx.rollback,
		profileStoreFunc(store.upsertProfile),
		ensurer.ensure,
		Profile{UserID: 42, Status: "active", DisplayName: "demo"},
	)
	if err != nil {
		t.Fatalf("runProfileUpsertTransaction() error = %v, want nil", err)
	}
	if tx.commits != 1 {
		t.Fatalf("runProfileUpsertTransaction() commits = %d, want 1", tx.commits)
	}
	if tx.rollbacks != 0 {
		t.Fatalf("runProfileUpsertTransaction() rollbacks = %d, want 0", tx.rollbacks)
	}
}
