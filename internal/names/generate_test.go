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
    "strings"

    "github.com/google/go-cmp/cmp"
    kerrors "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime/schema"
    "sigs.k8s.io/controller-runtime/pkg/client"

    "github.com/crossplane/crossplane-runtime/pkg/errors"
    "github.com/crossplane/crossplane-runtime/pkg/resource"
    "github.com/crossplane/crossplane-runtime/pkg/resource/fake"
    "github.com/crossplane/crossplane-runtime/pkg/test"
    "github.com/crossplane/crossplane/internal/xresource/unstructured/composite"
)

// fakeGVKComposed is a helper type to override GetObjectKind for testing purposes.
type fakeGVKComposed struct {
    *fake.Composed
}
// GetObjectKind returns a fixed GroupVersionKind.
func (f *fakeGVKComposed) GetObjectKind() schema.ObjectKind {
    return &metav1.TypeMeta{APIVersion: "test.group/v1", Kind: "TestKind"}
}
// fakeEmptyGVKComposed is a fake composed resource that returns an empty GroupVersionKind.
type fakeEmptyGVKComposed struct {
    *fake.Composed
}
// GetObjectKind returns an empty TypeMeta.
func (f *fakeEmptyGVKComposed) GetObjectKind() schema.ObjectKind {
    return &metav1.TypeMeta{}
}
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
        "ContextCanceled": {
            reason: "should return error if context is canceled",
            client: &test.MockClient{
                MockGet: func(ctx context.Context, key client.ObjectKey, obj client.Object) error {
                    return ctx.Err()
                },
            },
            args: args{
                ctx: func() context.Context {
                    ctx, cancel := context.WithCancel(context.Background())
                    cancel()
                    return ctx
                }(),
                cd: &fake.Composed{ObjectMeta: metav1.ObjectMeta{
                    GenerateName: "cool-resource-",
                }},
            },
            want: want{
                cd: &fake.Composed{ObjectMeta: metav1.ObjectMeta{
                    GenerateName: "cool-resource-",
                }},
                err: context.Canceled,
            },
        },
        "SuccessAfterMultipleConflicts": {
            reason: "Name is found after multiple conflicts",
            client: &test.MockClient{
                MockGet: func(ctx context.Context, key client.ObjectKey, obj client.Object) error {
                    if key.Name == "cool-resource-42" || key.Name == "cool-resource-43" {
                        return nil
                    }
                    return kerrors.NewNotFound(schema.GroupResource{Resource: "CoolResource"}, key.Name)
                },
            },
            args: args{
                ctx: context.Background(),
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
    "ContextDeadlineExceeded": {
        reason: "should return DeadlineExceeded if context deadline is exceeded",
        client: &test.MockClient{MockGet: func(ctx context.Context, key client.ObjectKey, obj client.Object) error {
            return context.DeadlineExceeded
        }},
        args: args{
            ctx: func() context.Context {
                ctx, cancel := context.WithTimeout(context.Background(), 0)
                defer cancel()
                return ctx
            }(),
            cd: &fake.Composed{ObjectMeta: metav1.ObjectMeta{
                GenerateName: "cool-resource-",
            }},
        },
        want: want{
            cd: &fake.Composed{ObjectMeta: metav1.ObjectMeta{
                GenerateName: "cool-resource-",
            }},
            err: context.DeadlineExceeded,
        },
    },
    "SetGroupVersionKind": {
        reason: "ensures the composite object gets the expected GroupVersionKind when generating a name",
        client: &test.MockClient{MockGet: func(ctx context.Context, key client.ObjectKey, obj client.Object) error {
            un, ok := obj.(*composite.Unstructured)
            if !ok {
                return errors.New("unexpected type")
            }
            expectedGVK := schema.GroupVersionKind{Group: "test.group", Version: "v1", Kind: "TestKind"}
            if got := un.GetObjectKind().GroupVersionKind(); got != expectedGVK {
                return errors.New("incorrect GVK")
            }
            return kerrors.NewNotFound(schema.GroupResource{Resource: "CoolResource"}, key.Name)
        }},
        args: args{
            cd: &fakeGVKComposed{&fake.Composed{ObjectMeta: metav1.ObjectMeta{
                GenerateName: "test-gvk-",
            }}},
        },
        want: want{
            cd: &fakeGVKComposed{&fake.Composed{ObjectMeta: metav1.ObjectMeta{
                GenerateName: "test-gvk-",
                Name:         "test-gvk-42",
            }}},
            err: nil,
        },
    },
    "ErrorAfterConflict": {
        reason: "Client returns error on retry after an initial conflict",
        client: &test.MockClient{MockGet: func(ctx context.Context, key client.ObjectKey, obj client.Object) error {
            if key.Name == "cool-resource-42" {
                // Simulate that the first generated name is taken.
                return nil
            }
            if key.Name == "cool-resource-43" {
                // Return a non-NotFound error on the next try.
                return errors.New("boom2")
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
            }},
            err: errors.New("boom2"),
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
func TestNameGeneratorFn(t *testing.T) {
    // This test verifies that the NameGeneratorFn type correctly implements the NameGenerator interface.
    // It provides a function that sets the name on a resource and ensures that the name is updated.
    fakeObj := &fake.Composed{ObjectMeta: metav1.ObjectMeta{GenerateName: "test-"}}
    // Define a NameGeneratorFn that sets a specific name.
    fn := NameGeneratorFn(func(ctx context.Context, cd resource.Object) error {
    cd.SetName("test-42")
    return nil
    })
    if err := fn.GenerateName(context.Background(), fakeObj); err != nil {
    t.Fatalf("unexpected error: %v", err)
    }
    if got, want := fakeObj.GetName(), "test-42"; got != want {
    t.Errorf("unexpected name; got %s, want %s", got, want)
    }
}
func TestNameGeneratorFnError(t *testing.T) {
    // TestNameGeneratorFnError verifies that a NameGeneratorFn returns an error when its function returns one.
    expectedErr := errors.New("test error")
    fn := NameGeneratorFn(func(ctx context.Context, cd resource.Object) error {
        return expectedErr
    })
    fakeObj := &fake.Composed{ObjectMeta: metav1.ObjectMeta{GenerateName: "test-"}}
    err := fn.GenerateName(context.Background(), fakeObj)
    if err != expectedErr {
        t.Errorf("expected error %v, got %v", expectedErr, err)
    }
}
type mockNameGenerator struct {
    last int
}

func (m *mockNameGenerator) GenerateName(prefix string) string {
    m.last++
    return prefix + strconv.Itoa(m.last)
}
// TestNewNameGenerator_Success tests that NewNameGenerator correctly generates a name by using the default SimpleNameGenerator.
func TestNewNameGenerator_Success(t *testing.T) {
    ctx := context.Background()
    // Create a composed object with GenerateName but no Name.
    fakeObj := &fake.Composed{ObjectMeta: metav1.ObjectMeta{GenerateName: "newtest-"}}

    // Create a mock client that always returns NotFound so that any generated name is available.
    c := &test.MockClient{
        MockGet: func(ctx context.Context, key client.ObjectKey, obj client.Object) error {
            return kerrors.NewNotFound(schema.GroupResource{Resource: "TestResource"}, key.Name)
        },
    }

    // Create a new name generator using the default SimpleNameGenerator.
    gen := NewNameGenerator(c)

    // Call GenerateName and verify no error is returned.
    if err := gen.GenerateName(ctx, fakeObj); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // Check that the generated name is non-empty and starts with the given prefix.
    got := fakeObj.GetName()
    if got == "" {
        t.Errorf("expected non-empty name, got empty")
    } else if !strings.HasPrefix(got, "newtest-") {
        t.Errorf("expected name to start with %q, got %q", "newtest-", got)
    }
}
// TestGenerateName_SkipAlreadyNamed_NoClientCall ensures that when an object already has a name,
// GenerateName does not attempt to generate one and therefore never calls the client's Get method.
func TestGenerateName_SkipAlreadyNamed_NoClientCall(t *testing.T) {
    called := false
    mockClient := &test.MockClient{
        MockGet: func(ctx context.Context, key client.ObjectKey, obj client.Object) error {
            called = true
            return nil
        },
    }
    gen := NewNameGenerator(mockClient)
    obj := &fake.Composed{ObjectMeta: metav1.ObjectMeta{
        Name:         "preset-name",
        GenerateName: "unused-",
    }}
    if err := gen.GenerateName(context.Background(), obj); err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if called {
        t.Errorf("client.Get was called unexpectedly")
    }
}
// TestGenerateName_EmptyGVK verifies that a resource with an empty GroupVersionKind
// still gets a name generated. This test increases coverage on the code path where
// the composite.Unstructured object's GVK is set using the value from cd.GetObjectKind().
func TestGenerateName_EmptyGVK(t *testing.T) {
    ctx := context.Background()
    // Create a resource that returns an empty GroupVersionKind.
    fakeObj := &fakeEmptyGVKComposed{
    Composed: &fake.Composed{ObjectMeta: metav1.ObjectMeta{
    GenerateName: "emptygvk-",

    }},
    }
    // Create a mock client that always returns a NotFound error so that any name is considered available.
    c := &test.MockClient{
    MockGet: func(ctx context.Context, key client.ObjectKey, obj client.Object) error {
    return kerrors.NewNotFound(schema.GroupResource{Resource: "TestResource"}, key.Name)
    },
    }
    // Use the mockNameGenerator to produce a predictable name.
    gen := &nameGenerator{reader: c, namer: &mockNameGenerator{last: 200}}
    err := gen.GenerateName(ctx, fakeObj)
    if err != nil {
    t.Fatalf("unexpected error: %v", err)
    }
    if got := fakeObj.GetName(); got == "" {
    t.Error("expected generated name to be non-empty, but got empty")
    } else if !strings.HasPrefix(got, "emptygvk-") {
    t.Errorf("expected generated name to start with %q, got %q", "emptygvk-", got)
    }
}