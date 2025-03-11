/*
Copyright 2023 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package names

import (
    "context"
    "strconv"
    "testing"

    "github.com/google/go-cmp/cmp"
    kerrors "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime/schema"
    "sigs.k8s.io/controller-runtime/pkg/client"

    "github.com/crossplane/crossplane-runtime/pkg/errors"
    "github.com/crossplane/crossplane-runtime/pkg/resource"
    "github.com/crossplane/crossplane-runtime/pkg/resource/fake"
    "github.com/crossplane/crossplane-runtime/pkg/test"
)

func TestGenerateName(t *testing.T) {
    errBoom := errors.New("boom")

    type args struct {
    ctx context.Context
    cd  resource.Composed
    }
    type want struct {
    cd  resource.Composed
    err error
    }
    cases := map[string]struct {
    reason string
    client client.Client
    args
    want
    }{
    "SkipGenerateNamedResources": {
    reason: "We should not try naming a resource that already have a name",
    // We must be returning early, or else we'd hit this error.
    client: &test.MockClient{MockCreate: test.NewMockCreateFn(errBoom)},
    args: args{
    cd: &fake.Composed{ObjectMeta: metav1.ObjectMeta{
        Name: "already-has-a-cool-name",
    }},
    },
    want: want{
    cd: &fake.Composed{ObjectMeta: metav1.ObjectMeta{
        Name: "already-has-a-cool-name",
    }},
    err: nil,
    },
    },
    "SkipGenerateNameForResourcesWithoutGenerateName": {
    reason: "We should not try to name resources that don't have a generate name (though that should never happen)",
    // We must be returning early, or else we'd hit this error.
    client: &test.MockClient{MockCreate: test.NewMockCreateFn(errBoom)},
    args: args{
    cd: &fake.Composed{}, // Conspicously missing a generate name.
    },
    want: want{
    cd:  &fake.Composed{},
    err: nil,
    },
    },
    "NameGeneratorClientError": {
    reason: "Client error finding a free name for a composed resource",
    client: &test.MockClient{MockGet: test.NewMockGetFn(errBoom)},
    args: args{
    cd: &fake.Composed{ObjectMeta: metav1.ObjectMeta{
        GenerateName: "cool-resource-",
    }},
    },
    want: want{
    cd: &fake.Composed{ObjectMeta: metav1.ObjectMeta{
        GenerateName: "cool-resource-",
    }},
    err: errBoom,
    },
    },
    "Success": {
    reason: "Name is found on first try",
    client: &test.MockClient{MockGet: test.NewMockGetFn(kerrors.NewNotFound(schema.GroupResource{Resource: "CoolResource"}, "cool-resource-42"))},
    args: args{
    cd: &fake.Composed{ObjectMeta: metav1.ObjectMeta{
        GenerateName: "cool-resource-",
    }},
    },
    want: want{
    cd: &fake.Composed{ObjectMeta: metav1.ObjectMeta{
        GenerateName: "cool-resource-",
        Name:         "cool-resource-42",
    }},
    },
    },
    "SuccessAfterConflict": {
    reason: "Name is found on second try",
    client: &test.MockClient{MockGet: func(_ context.Context, key client.ObjectKey, _ client.Object) error {
    if key.Name == "cool-resource-42" {
        return nil
    }
    return kerrors.NewNotFound(schema.GroupResource{Resource: "CoolResource"}, key.Name)
    }},
    args: args{
    cd: &fake.Composed{ObjectMeta: metav1.ObjectMeta{
        GenerateName: "cool-resource-",
    }},
    },
    want: want{
    cd: &fake.Composed{ObjectMeta: metav1.ObjectMeta{
        GenerateName: "cool-resource-",
        Name:         "cool-resource-43",
    }},
        },
    },
    "SuccessAfterTwoConflicts": {
        reason: "Name is found after two conflicts",
        client: &test.MockClient{MockGet: func(ctx context.Context, key client.ObjectKey, obj client.Object) error {
            if key.Name == "cool-resource-42" || key.Name == "cool-resource-43" {
                return nil
            }
            return kerrors.NewNotFound(schema.GroupResource{Resource: "CoolResource"}, key.Name)
        }},
        args: args{
            cd: &fake.Composed{ObjectMeta: metav1.ObjectMeta{
                GenerateName: "cool-resource-",
            }},
        },
        want: want{
            cd: &fake.Composed{ObjectMeta: metav1.ObjectMeta{
                GenerateName: "cool-resource-",
                Name:         "cool-resource-44",
            }},
            err: nil,
        },
    },
    "AlwaysConflict": {
    reason: "Name cannot be found",
    client: &test.MockClient{MockGet: test.NewMockGetFn(nil)},
    args: args{
    cd: &fake.Composed{ObjectMeta: metav1.ObjectMeta{
        GenerateName: "cool-resource-",
    }},
    },
    want: want{
    cd: &fake.Composed{ObjectMeta: metav1.ObjectMeta{
        GenerateName: "cool-resource-",
    }},
    err: errors.New(errGenerateName),
    },
    },
    }
    for name, tc := range cases {
    t.Run(name, func(t *testing.T) {
    r := &nameGenerator{reader: tc.client, namer: &mockNameGenerator{last: 41}}
    err := r.GenerateName(tc.args.ctx, tc.args.cd)
    if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
    t.Errorf("\n%s\nDryRunRender(...): -want, +got:\n%s", tc.reason, diff)
    }
    if diff := cmp.Diff(tc.want.cd, tc.args.cd); diff != "" {
    t.Errorf("\n%s\nDryRunRender(...): -want, +got:\n%s", tc.reason, diff)
    }
    })
    }
}
// TestNameGeneratorFn validates that a NameGeneratorFn calls the underlying function and returns the expected result.
func TestNameGeneratorFn(t *testing.T) {
    ctx := context.Background()

    t.Run("returns expected error and sets name", func(t *testing.T) {
        // Test a NameGeneratorFn that sets a name and returns an error.
        expectedErr := errors.New("custom error")
        cd := &fake.Composed{}
        // Define a NameGeneratorFn that sets the name and returns an error.
        fn := NameGeneratorFn(func(ctx context.Context, cd resource.Object) error {
            cd.SetName("generated-name")
            return expectedErr
        })

        err := fn.GenerateName(ctx, cd)
        if diff := cmp.Diff(expectedErr, err, test.EquateErrors()); diff != "" {
            t.Errorf("unexpected error (-want, +got):\n%s", diff)
        }
        if n := cd.GetName(); n != "generated-name" {
            t.Errorf("expected name generated-name, got %q", n)
        }
    })

    t.Run("returns nil when no error and sets name", func(t *testing.T) {
        // Test a NameGeneratorFn that sets a name and returns no error.
        cd := &fake.Composed{}
        fn := NameGeneratorFn(func(ctx context.Context, cd resource.Object) error {
            cd.SetName("non-error-name")
            return nil
        })

        err := fn.GenerateName(ctx, cd)
        if err != nil {
            t.Errorf("expected no error, but got: %v", err)
        }
        if n := cd.GetName(); n != "non-error-name" {
            t.Errorf("expected name non-error-name, got %q", n)
        }
    })
}
// TestGenerateNameEmptyGVK tests that even if the resource's GroupVersionKind is empty,
func TestGenerateNameEmptyGVK(t *testing.T) {
    ctx := context.Background()
    // Create a fake Composed resource with a non-empty GenerateName and explicitly set an empty GVK.
    cd := &fake.Composed{ObjectMeta: metav1.ObjectMeta{
    GenerateName: "test-empty-gvk-",
    }}
    // Use a MockClient that returns a NotFound error so that the generated name is accepted.
    client := &test.MockClient{
    MockGet: test.NewMockGetFn(kerrors.NewNotFound(schema.GroupResource{Resource: "TestResource"}, "test-empty-gvk-1")),
    }
    // Use our mock name generator starting from 0.
    gen := &nameGenerator{reader: client, namer: &mockNameGenerator{last: 0}}

    if err := gen.GenerateName(ctx, cd); err != nil {
    t.Errorf("expected no error, got: %v", err)
    }

    // Our mockNameGenerator increments last and returns prefix+"1" on first call.
    if got, want := cd.GetName(), "test-empty-gvk-1"; got != want {
    t.Errorf("expected name %q, got %q", want, got)
    }
}

// TestNewNameGenerator tests that NewNameGenerator returns a valid NameGenerator
func TestNewNameGenerator(t *testing.T) {
    ctx := context.Background()
    // Create a MockClient with no side effects.
    client := &test.MockClient{}
    // Create a NameGenerator via the public API.
    gen := NewNameGenerator(client)
    if gen == nil {
    t.Fatal("expected non-nil name generator")
    }

    // Create a resource that already has a name.
    cd := &fake.Composed{ObjectMeta: metav1.ObjectMeta{
    Name:         "already-exists",
    GenerateName: "should-not-be-used-",
    }}
    // Call GenerateName. Since cd already has a non-empty Name,
    // we expect GenerateName to do nothing.
    if err := gen.GenerateName(ctx, cd); err != nil {
    t.Errorf("expected no error, got: %v", err)
    }
    if got, want := cd.GetName(), "already-exists"; got != want {
    t.Errorf("expected name %q, got %q", want, got)
    }
}
// TestGenerateNameNilResource verifies that passing a nil resource causes a panic.
func TestGenerateNameNilResource(t *testing.T) {
    ctx := context.Background()
    // Use a dummy MockClient since the nil resource will trigger a panic before client.Get is called.
    client := &test.MockClient{}
    gen := &nameGenerator{reader: client, namer: &mockNameGenerator{last: 0}}

    defer func() {
    if r := recover(); r == nil {
    t.Errorf("expected panic when passing nil resource, but no panic occurred")
    }
    }()

    // Pass nil resource; this should cause a panic since cd.GetName() is invoked.
    gen.GenerateName(ctx, nil)
    t.Errorf("should have panicked when nil resource was passed")
}

// TestGenerateNameGetErrorAfterConflict verifies that if the first call returns a conflict (i.e. a name is already in use)
// and the second call returns a non-NotFound error, GenerateName will return that error immediately.
func TestGenerateNameGetErrorAfterConflict(t *testing.T) {
    ctx := context.Background()
    callCount := 0
    expectedErr := errors.New("boom after conflict")

    // Set up a MockClient that returns nil (indicating the name is taken) on the first call,
    // and returns a custom error on the second call.
    client := &test.MockClient{
    MockGet: func(ctx context.Context, key client.ObjectKey, obj client.Object) error {
    callCount++
    if callCount == 1 {
    // Simulate that the name is already taken (client.Get returns nil, not a NotFound error).
    return nil
    }
    // On subsequent calls, return a non-NotFound error.
    return expectedErr
    },
    }

    cd := &fake.Composed{ObjectMeta: metav1.ObjectMeta{
    GenerateName: "test-resource-",
    }}

    // Initialize our nameGenerator with the mock namer starting at 0.
    gen := &nameGenerator{reader: client, namer: &mockNameGenerator{last: 0}}

    err := gen.GenerateName(ctx, cd)
    if err != expectedErr {
    t.Errorf("expected error %v, got %v", expectedErr, err)
    }

    // We expect that exactly 2 attempts were made: one conflict and one error.
    if callCount != 2 {
    t.Errorf("expected 2 attempts, got %d", callCount)
    }

    // Since the generator aborted on an error, no name should be set on the resource.
    if cd.GetName() != "" {
    t.Errorf("expected no generated name due to error, but got %q", cd.GetName())
    }
}
// TestGenerateName_NoClientCallIfNameSet verifies that GenerateName does not call client.Get when the resource already has a name.
func TestGenerateName_NoClientCallIfNameSet(t *testing.T) {
    ctx := context.Background()
    cd := &fake.Composed{ObjectMeta: metav1.ObjectMeta{
        Name:         "pre-set-name",
        GenerateName: "should-not-be-used-",
    }}

    // Define a custom client that will fail the test if Get is ever called.

    nc := &noCallClient{}
    // We pass a mock name generator as well; its starting value shouldn't matter since no generation is expected.
    gen := &nameGenerator{reader: nc, namer: &mockNameGenerator{last: 100}}
    err := gen.GenerateName(ctx, cd)
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    // The name should remain unchanged since the resource already had a name.
    if got, want := cd.GetName(), "pre-set-name"; got != want {
        t.Errorf("expected name %q, got %q", want, got)
    }
}
type mockNameGenerator struct {
    last int
}

func (m *mockNameGenerator) GenerateName(prefix string) string {
    m.last++
    return prefix + strconv.Itoa(m.last)
}
// noCallClient is a package-level mock client that panics if Get is called.
type noCallClient struct{}

func (c *noCallClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
    panic("client.Get should not be called when resource already has a name")
}

func (c *noCallClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
    return nil
}

func (c *noCallClient) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
    return nil
}

func (c *noCallClient) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
    return nil
}

func (c *noCallClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
    return nil
}

func (c *noCallClient) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
    return nil
}

func (c *noCallClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
    return nil
}

func (c *noCallClient) Status() client.StatusWriter { 
    return nil 
}