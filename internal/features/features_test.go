package features

import (
    "testing"

    "github.com/crossplane/crossplane-runtime/pkg/feature"
    "encoding/json"
    "fmt"
)

// TestAlphaFeatureFlags verifies that the alpha feature flags have the expected string values.
func TestAlphaFeatureFlags(t *testing.T) {
    tests := []struct {
        flag feature.Flag
        want string
    }{
        {EnableAlphaRealtimeCompositions, "EnableAlphaRealtimeCompositions"},
        {EnableAlphaDependencyVersionUpgrades, "EnableAlphaDependencyVersionUpgrades"},
        {EnableAlphaSignatureVerification, "EnableAlphaSignatureVerification"},
    }
    for _, tt := range tests {
        if string(tt.flag) != tt.want {
            t.Errorf("expected flag %q, got %q", tt.want, tt.flag)
        }
    }
}

// TestBetaFeatureFlags verifies that the beta feature flags have the expected string values.
func TestBetaFeatureFlags(t *testing.T) {
    tests := []struct {
        flag feature.Flag
        want string
    }{
        {EnableBetaDeploymentRuntimeConfigs, "EnableBetaDeploymentRuntimeConfigs"},
        {EnableBetaUsages, "EnableBetaUsages"},
        {EnableBetaClaimSSA, "EnableBetaClaimSSA"},
    }
    for _, tt := range tests {
        if string(tt.flag) != tt.want {
            t.Errorf("expected flag %q, got %q", tt.want, tt.flag)
        }
    }
}

// TestUniqueFlags checks that all feature flags have unique string values.
func TestUniqueFlags(t *testing.T) {
    flagMap := map[string]bool{
        string(EnableAlphaRealtimeCompositions):    true,
        string(EnableAlphaDependencyVersionUpgrades): true,
        string(EnableAlphaSignatureVerification):     true,
        string(EnableBetaDeploymentRuntimeConfigs):   true,
        string(EnableBetaUsages):                     true,
        string(EnableBetaClaimSSA):                   true,
    }
    if len(flagMap) != 6 {
        t.Errorf("expected 6 unique flag values, got %d", len(flagMap))
    }
}

// TestFlagBehavior is a dummy test that shows usage of a feature flag check.
func TestFlagBehavior(t *testing.T) {
    // Although our flags are constants, we simulate a check.
    if EnableAlphaRealtimeCompositions != "EnableAlphaRealtimeCompositions" {
        t.Errorf("Flag value mismatch: expected %v", "EnableAlphaRealtimeCompositions")
    }
}
// TestJSONMarshalling verifies that each feature flag can be marshalled to JSON and then unmarshalled back.
func TestJSONMarshalling(t *testing.T) {
    tests := []feature.Flag{
        EnableAlphaRealtimeCompositions,
        EnableAlphaDependencyVersionUpgrades,
        EnableAlphaSignatureVerification,
        EnableBetaDeploymentRuntimeConfigs,
        EnableBetaUsages,
        EnableBetaClaimSSA,
    }
    for _, flag := range tests {
        jsonData, err := json.Marshal(flag)
        if err != nil {
            t.Errorf("failed to marshal flag %v: %v", flag, err)
        }
        var unmarshalled string
        if err := json.Unmarshal(jsonData, &unmarshalled); err != nil {
            t.Errorf("failed to unmarshal json data %v: %v", string(jsonData), err)
        }
        if unmarshalled != string(flag) {
            t.Errorf("expected unmarshalled flag %q, got %q", flag, unmarshalled)
        }
    }
}

// TestFlagInStructJSONMarshalling verifies that feature flags embedded in a struct are correctly encoded to JSON.
func TestFlagInStructJSONMarshalling(t *testing.T) {
    type testStruct struct {
        Feature feature.Flag `json:"feature"`
        Desc    string       `json:"desc"`
    }
    s := testStruct{
        Feature: EnableAlphaSignatureVerification,
        Desc:    "test description",
    }
    b, err := json.Marshal(s)
    if err != nil {
        t.Errorf("failed to marshal test struct: %v", err)
    }
    want := `{"feature":"EnableAlphaSignatureVerification","desc":"test description"}`
    if string(b) != want {
        t.Errorf("expected JSON %q, got %q", want, string(b))
    }
}
// TestInvalidJSONUnmarshalling verifies that invalid JSON input for feature.Flag returns an error.
func TestInvalidJSONUnmarshalling(t *testing.T) {
    var flag feature.Flag
    invalidJSON := []byte(`123`) // invalid JSON: non-string value
    err := json.Unmarshal(invalidJSON, &flag)
    if err == nil {
        t.Error("expected error when unmarshalling invalid JSON into feature.Flag, got nil")
    }
}

// TestFlagPointerMarshalling verifies that a feature.Flag embedded as a pointer in a struct is correctly marshalled and unmarshalled.
func TestFlagPointerMarshalling(t *testing.T) {
    type pointerStruct struct {
        Feature *feature.Flag `json:"feature,omitempty"`
        Desc    string        `json:"desc"`
    }
    flag := EnableBetaUsages
    s := pointerStruct{
        Feature: &flag,
        Desc:    "pointer test",
    }
    data, err := json.Marshal(s)
    if err != nil {
        t.Errorf("failed to marshal pointer struct: %v", err)
    }
    expected := `{"feature":"EnableBetaUsages","desc":"pointer test"}`
    if string(data) != expected {
        t.Errorf("expected JSON %q, got %q", expected, string(data))
    }
    var s2 pointerStruct
    if err := json.Unmarshal(data, &s2); err != nil {
        t.Errorf("failed to unmarshal JSON into struct: %v", err)
    }
    if s2.Feature == nil || *s2.Feature != EnableBetaUsages {
        t.Errorf("expected Feature %q after unmarshal, got %v", "EnableBetaUsages", s2.Feature)
    }
}

// TestEmptyFlagMarshalling verifies that an empty feature.Flag can be marshalled and unmarshalled correctly.
func TestEmptyFlagMarshalling(t *testing.T) {
    emptyFlag := feature.Flag("")
    data, err := json.Marshal(emptyFlag)
    if err != nil {
        t.Errorf("failed to marshal empty flag: %v", err)
    }
    if string(data) != `""` {
        t.Errorf("expected empty JSON string, got %q", string(data))
    }
    var unmarshalled feature.Flag
    if err := json.Unmarshal(data, &unmarshalled); err != nil {
        t.Errorf("failed to unmarshal empty flag: %v", err)
    }
    if unmarshalled != emptyFlag {
        t.Errorf("expected empty flag %q after unmarshal, got %q", emptyFlag, unmarshalled)
    }
}

// TestFlagEquality verifies that feature.Flag comparisons work as expected.
func TestFlagEquality(t *testing.T) {
    // Compare constant flag with expected string literal.
    if EnableAlphaRealtimeCompositions != feature.Flag("EnableAlphaRealtimeCompositions") {
        t.Errorf("flag equality check failed for EnableAlphaRealtimeCompositions")
    }
    if EnableBetaClaimSSA != feature.Flag("EnableBetaClaimSSA") {
        t.Errorf("flag equality check failed for EnableBetaClaimSSA")
    }
}
// TestNilFlagPointerMarshalling verifies that if a pointer feature.Flag is nil, it is omitted in JSON marshalling.
func TestNilFlagPointerMarshalling(t *testing.T) {
    type pointerStruct struct {
        Feature *feature.Flag `json:"feature,omitempty"`
        Desc    string        `json:"desc"`
    }
    s := pointerStruct{
        Feature: nil,
        Desc:    "nil pointer test",
    }
    data, err := json.Marshal(s)
    if err != nil {
        t.Errorf("failed to marshal struct with nil feature flag: %v", err)
    }
    expected := `{"desc":"nil pointer test"}`
    if string(data) != expected {
        t.Errorf("expected JSON %q, got %q", expected, string(data))
    }

    var s2 pointerStruct
    if err := json.Unmarshal(data, &s2); err != nil {
        t.Errorf("failed to unmarshal JSON into struct with nil feature flag: %v", err)
    }
    if s2.Feature != nil {
        t.Errorf("expected nil feature flag after unmarshal, got %v", *s2.Feature)
    }
}

// TestDirectJSONUnmarshalling verifies that unmarshalling a valid JSON string directly into a feature.Flag works correctly.
func TestDirectJSONUnmarshalling(t *testing.T) {
    var f feature.Flag
    jsonData := []byte(`"CustomFlagValue"`)
    if err := json.Unmarshal(jsonData, &f); err != nil {
        t.Errorf("failed to unmarshal JSON into feature.Flag: %v", err)
    }
    if f != feature.Flag("CustomFlagValue") {
        t.Errorf("expected feature flag %q, got %q", "CustomFlagValue", f)
    }
}
// TestFlagPointerJSONNullUnmarshal verifies that unmarshalling a JSON object with a null feature pointer sets the feature field to nil.
func TestFlagPointerJSONNullUnmarshal(t *testing.T) {
    type pointerStruct struct {
        Feature *feature.Flag `json:"feature,omitempty"`
        Desc    string        `json:"desc"`
    }
    // JSON with feature explicitly set to null.
    jsonStr := `{"feature": null, "desc": "null feature test"}`
    var s pointerStruct
    if err := json.Unmarshal([]byte(jsonStr), &s); err != nil {
        t.Errorf("failed to unmarshal JSON with null feature: %v", err)
    }
    if s.Feature != nil {
        t.Errorf("expected nil feature flag after unmarshalling JSON null, got: %v", *s.Feature)
    }
}

// TestSliceOfFlagsJSONMarshalling verifies that a slice of feature.Flags is marshalled to a valid JSON array of strings.
func TestSliceOfFlagsJSONMarshalling(t *testing.T) {
    flags := []feature.Flag{
        EnableAlphaRealtimeCompositions,
        EnableAlphaDependencyVersionUpgrades,
        EnableAlphaSignatureVerification,
        EnableBetaDeploymentRuntimeConfigs,
        EnableBetaUsages,
        EnableBetaClaimSSA,
    }
    jsonData, err := json.Marshal(flags)
    if err != nil {
        t.Fatalf("failed to marshal slice of flags: %v", err)
    }
    // Build expected JSON array.
    expected := `["EnableAlphaRealtimeCompositions","EnableAlphaDependencyVersionUpgrades","EnableAlphaSignatureVerification",` +
        `"EnableBetaDeploymentRuntimeConfigs","EnableBetaUsages","EnableBetaClaimSSA"]`
    if string(jsonData) != expected {
        t.Errorf("expected JSON %q, got %q", expected, string(jsonData))
    }
}
// TestNestedStructJSONMarshalling verifies that a nested struct containing feature.Flags in slices, maps, and nested structs
// is correctly encoded to JSON and then decoded with the values remaining unchanged.
func TestNestedStructJSONMarshalling(t *testing.T) {
    type Nested struct {
    Title string      `json:"title"`
    Flag  feature.Flag `json:"flag"`
    }
    type ComplexStruct struct {
    Flags   []feature.Flag          `json:"flags"`
    FlagMap map[string]feature.Flag `json:"flag_map"`
    Nested  Nested                  `json:"nested"`
    }

    cs := ComplexStruct{
    Flags: []feature.Flag{
    EnableAlphaRealtimeCompositions,
    EnableBetaUsages,
    },
    FlagMap: map[string]feature.Flag{
    "alpha": EnableAlphaSignatureVerification,
    "beta":  EnableBetaClaimSSA,
    },
    Nested: Nested{
    Title: "Nested Test",
    Flag:  EnableBetaDeploymentRuntimeConfigs,
    },
    }

    data, err := json.Marshal(cs)
    if err != nil {
    t.Fatalf("failed to marshal complex struct: %v", err)
    }

    var cs2 ComplexStruct
    if err := json.Unmarshal(data, &cs2); err != nil {
    t.Fatalf("failed to unmarshal complex struct: %v", err)
    }

    if len(cs2.Flags) != 2 {
    t.Errorf("expected 2 flags in slice, got %d", len(cs2.Flags))
    } else {
    if cs2.Flags[0] != EnableAlphaRealtimeCompositions {
    t.Errorf("expected first flag %q, got %q", EnableAlphaRealtimeCompositions, cs2.Flags[0])
    }
    if cs2.Flags[1] != EnableBetaUsages {
    t.Errorf("expected second flag %q, got %q", EnableBetaUsages, cs2.Flags[1])
    }
    }

    if len(cs2.FlagMap) != 2 ||
    cs2.FlagMap["alpha"] != EnableAlphaSignatureVerification ||
    cs2.FlagMap["beta"] != EnableBetaClaimSSA {
    t.Errorf("unexpected flag map contents: %v", cs2.FlagMap)
    }

    if cs2.Nested.Title != "Nested Test" || cs2.Nested.Flag != EnableBetaDeploymentRuntimeConfigs {
    t.Errorf("unexpected nested struct: %+v", cs2.Nested)
    }
}
// TestFlagCaseSensitivity verifies that flag comparisons are case-sensitive.
func TestFlagCaseSensitivity(t *testing.T) {
    // Create a lower-case version of the constant.
    lower := feature.Flag("enablealpharealtimecompositions")
    if EnableAlphaRealtimeCompositions == lower {
        t.Errorf("expected flag comparison to be case-sensitive and not equal; got %q equal to %q", EnableAlphaRealtimeCompositions, lower)
    }
}

// TestSpecialCharacterFlag verifies that JSON marshalling and unmarshalling work correctly for a flag containing special characters.
func TestSpecialCharacterFlag(t *testing.T) {
    // Create a custom flag with quotes and newline characters.
    special := feature.Flag("Test\"Special\nFlag")
    jsonData, err := json.Marshal(special)
    if err != nil {
        t.Fatalf("failed to marshal special flag: %v", err)
    }
    var unmarshalled feature.Flag
    if err := json.Unmarshal(jsonData, &unmarshalled); err != nil {
        t.Fatalf("failed to unmarshal special flag: %v", err)
    }
    if special != unmarshalled {
        t.Errorf("expected special flag %q, got %q", special, unmarshalled)
    }
}

// TestFlagFormat verifies that the string representation obtained via fmt.Sprintf matches the flag value.
func TestFlagFormat(t *testing.T) {
    flagVal := EnableBetaClaimSSA
    formatted := fmt.Sprintf("%s", flagVal)
    if formatted != string(flagVal) {
        t.Errorf("expected formatted flag %q to equal %q", formatted, flagVal)
    }
}
// TestFlagAsMapKey verifies that feature.Flag constants can be used as map keys and retrieved correctly.
func TestFlagAsMapKey(t *testing.T) {
    m := map[feature.Flag]int{
        EnableAlphaRealtimeCompositions: 1,
        EnableBetaClaimSSA:              2,
    }

    if m[EnableAlphaRealtimeCompositions] != 1 {
        t.Errorf("expected value 1 for flag %q, got %d", EnableAlphaRealtimeCompositions, m[EnableAlphaRealtimeCompositions])
    }

    if m[EnableBetaClaimSSA] != 2 {
        t.Errorf("expected value 2 for flag %q, got %d", EnableBetaClaimSSA, m[EnableBetaClaimSSA])
    }
}

// TestFlagMapKeyMarshalling verifies that a map with feature.Flag keys is correctly marshalled to JSON,
// and that the resulting JSON keys match the string representations of those feature flags.
func TestFlagMapKeyMarshalling(t *testing.T) {
    m := map[feature.Flag]string{
        EnableAlphaDependencyVersionUpgrades: "alpha-dep",
        EnableBetaUsages:                     "beta-usages",
    }

    jsonData, err := json.Marshal(m)
    if err != nil {
        t.Fatalf("failed to marshal map with feature.Flag keys: %v", err)
    }

    // Unmarshal the JSON data into a generic map[string]string because JSON object keys are strings.
    var m2 map[string]string
    if err := json.Unmarshal(jsonData, &m2); err != nil {
        t.Fatalf("failed to unmarshal JSON data into map[string]string: %v", err)
    }

    if v, ok := m2[string(EnableAlphaDependencyVersionUpgrades)]; !ok || v != "alpha-dep" {
        t.Errorf("expected key %q with value %q, got %q", EnableAlphaDependencyVersionUpgrades, "alpha-dep", v)
    }

    if v, ok := m2[string(EnableBetaUsages)]; !ok || v != "beta-usages" {
        t.Errorf("expected key %q with value %q, got %q", EnableBetaUsages, "beta-usages", v)
    }
}
// TestFlagRoundTripStringConversion verifies that converting a feature.Flag to a string and back retains equality.
func TestFlagRoundTripStringConversion(t *testing.T) {
    original := feature.Flag("TestRoundTrip")
    str := string(original)
    roundTrip := feature.Flag(str)
    if original != roundTrip {
        t.Errorf("expected round trip conversion to yield same flag, got original %q and roundTrip %q", original, roundTrip)
    }
}

// TestRepeatedJSONMarshalling verifies that repeated JSON marshalling and unmarshalling of a feature.Flag remains consistent.
func TestRepeatedJSONMarshalling(t *testing.T) {
    flagVal := EnableAlphaRealtimeCompositions
    for i := 0; i < 10; i++ {
        data, err := json.Marshal(flagVal)
        if err != nil {
            t.Fatalf("iteration %d: failed to marshal flag: %v", i, err)
        }
        var unmarshalled feature.Flag
        err = json.Unmarshal(data, &unmarshalled)
        if err != nil {
            t.Fatalf("iteration %d: failed to unmarshal flag: %v", i, err)
        }
        if flagVal != unmarshalled {
            t.Errorf("iteration %d: expected marshalled value %q, got %q", i, flagVal, unmarshalled)
        }
    }
}