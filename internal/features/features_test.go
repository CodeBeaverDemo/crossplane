package features

import (
    "testing"
    "strings"
    "sort"

    "github.com/crossplane/crossplane-runtime/pkg/feature"
    "encoding/json"
    "fmt"
)

// TestAlphaFlags tests that each alpha feature flag has the expected value.
func TestAlphaFlags(t *testing.T) {
    // Define test cases for alpha feature flags.
    testCases := []struct {
    flag     feature.Flag
    expected string
    }{
    {EnableAlphaRealtimeCompositions, "EnableAlphaRealtimeCompositions"},
    {EnableAlphaDependencyVersionUpgrades, "EnableAlphaDependencyVersionUpgrades"},
    {EnableAlphaSignatureVerification, "EnableAlphaSignatureVerification"},
    }
    for _, tc := range testCases {
    if tc.flag != feature.Flag(tc.expected) {
    t.Errorf("Expected alpha feature flag %q, got %q", tc.expected, tc.flag)
    }
    }
}

// TestBetaFlags tests that each beta feature flag has the expected string value.
func TestBetaFlags(t *testing.T) {
    // Define test cases for beta feature flags.
    testCases := []struct {
    flag     feature.Flag
    expected string
    }{
    {EnableBetaDeploymentRuntimeConfigs, "EnableBetaDeploymentRuntimeConfigs"},
    {EnableBetaUsages, "EnableBetaUsages"},
    {EnableBetaClaimSSA, "EnableBetaClaimSSA"},
    }
    for _, tc := range testCases {
    if tc.flag != feature.Flag(tc.expected) {
    t.Errorf("Expected beta feature flag %q, got %q", tc.expected, tc.flag)
    }
    }
}
// TestFlagMapKeys tests that feature flags can be used as map keys
// and that their string values contain meaningful substrings.
func TestFlagMapKeys(t *testing.T) {
    // Create a map with feature.Flag as keys and an expected substring value.
    flags := map[feature.Flag]string{
        EnableAlphaRealtimeCompositions:     "RealtimeCompositions",
        EnableAlphaDependencyVersionUpgrades:  "DependencyVersionUpgrades",
        EnableAlphaSignatureVerification:      "SignatureVerification",
        EnableBetaDeploymentRuntimeConfigs:    "DeploymentRuntimeConfigs",
        EnableBetaUsages:                      "Usages",
        EnableBetaClaimSSA:                    "ClaimSSA",
    }

    // Verify that the keys in the map have the expected substring in their string value.
    for flag, expectedSubstring := range flags {
        if !strings.Contains(string(flag), expectedSubstring) {
            t.Errorf("Expected flag %q to contain substring %q", flag, expectedSubstring)
        }
    }
}

// TestUnknownFlag tests that an unknown feature flag behaves like a normal string.
func TestUnknownFlag(t *testing.T) {
    unknown := feature.Flag("UnknownFeature")
    if string(unknown) != "UnknownFeature" {
        t.Errorf("Expected unknown feature flag to be %q, got %q", "UnknownFeature", unknown)
    }
}
// TestEmptyFlag tests that an empty feature flag returns an empty string.
func TestEmptyFlag(t *testing.T) {
    empty := feature.Flag("")
    if empty != "" {
        t.Errorf("Expected empty flag to be empty, but got %q", empty)
    }
}

// TestFlagEquality tests that feature flags with the same values are equal,
// and flags with different values are not equal.
func TestFlagEquality(t *testing.T) {
    // Create two flags with identical values.
    flag1 := feature.Flag("TestFlag")
    flag2 := feature.Flag("TestFlag")

    // Create a flag with a different value.
    flag3 := feature.Flag("DifferentTestFlag")

    if flag1 != flag2 {
        t.Errorf("Expected flags %q and %q to be equal", flag1, flag2)
    }
    if flag1 == flag3 {
        t.Errorf("Expected flags %q and %q to be different", flag1, flag3)
    }
}

// TestFlagConcatenation tests that feature.Flag can be concatenated with strings.
func TestFlagConcatenation(t *testing.T) {
    flag := feature.Flag("ConcatTest")
    result := "Prefix-" + string(flag) + "-Suffix"
    expected := "Prefix-ConcatTest-Suffix"
    if result != expected {
        t.Errorf("Expected concatenated result %q, got %q", expected, result)
    }
}
// TestFlagJSONMarshalling tests that a feature.Flag can be marshaled and unmarshaled correctly.
func TestFlagJSONMarshalling(t *testing.T) {
    type FlagHolder struct {
        Flag feature.Flag `json:"flag"`
    }

    original := FlagHolder{Flag: EnableAlphaSignatureVerification}
    data, err := json.Marshal(original)
    if err != nil {
        t.Fatalf("Failed to marshal: %v", err)
    }

    var decoded FlagHolder
    if err := json.Unmarshal(data, &decoded); err != nil {
        t.Fatalf("Failed to unmarshal: %v", err)
    }

    if original.Flag != decoded.Flag {
        t.Errorf("Expected flag %q after JSON roundtrip, got %q", original.Flag, decoded.Flag)
    }
}

// TestFlagSpecialCharacters tests that feature.Flag correctly handles strings with special characters.
func TestFlagSpecialCharacters(t *testing.T) {
    special := feature.Flag("Special_!@#$%^&*()_+-=[]{}|;':,./<>?")
    if string(special) != "Special_!@#$%^&*()_+-=[]{}|;':,./<>?" {
        t.Errorf("Expected special flag to be %q, got %q", "Special_!@#$%^&*()_+-=[]{}|;':,./<>?", special)
    }
}
// TestFlagDirectJSONMarshalling tests that a feature.Flag is marshalled as a JSON string.
func TestFlagDirectJSONMarshalling(t *testing.T) {
    flag := feature.Flag("DirectTest")
    data, err := json.Marshal(flag)
    if err != nil {
        t.Fatalf("Failed to marshal flag: %v", err)
    }
    expected := "\"DirectTest\""
    if string(data) != expected {
        t.Errorf("Expected JSON output %q, got %q", expected, string(data))
    }
}

// TestFlagDirectJSONUnmarshalling tests the unmarshalling of a JSON string directly into a feature.Flag.
func TestFlagDirectJSONUnmarshalling(t *testing.T) {
    jsonData := []byte("\"UnmarshalledFlag\"")
    var flag feature.Flag
    if err := json.Unmarshal(jsonData, &flag); err != nil {
        t.Fatalf("Failed to unmarshal JSON into flag: %v", err)
    }
    expected := "UnmarshalledFlag"
    if string(flag) != expected {
        t.Errorf("Expected flag %q, got %q", expected, flag)
    }
}
// TestEmptyFlagJSONMarshalling tests that an empty feature.Flag is marshaled correctly as a JSON empty string.
func TestEmptyFlagJSONMarshalling(t *testing.T) {
    emptyFlag := feature.Flag("")
    data, err := json.Marshal(emptyFlag)
    if err != nil {
        t.Fatalf("Failed to marshal empty flag: %v", err)
    }
    expected := "\"\""
    if string(data) != expected {
        t.Errorf("Expected JSON output %q for empty flag, got %q", expected, string(data))
    }
}

// TestFlagUnmarshalInvalidJSON tests that unmarshalling an invalid JSON type into a feature.Flag returns an error.
func TestFlagUnmarshalInvalidJSON(t *testing.T) {
    // Here we use a JSON number which is invalid for unmarshalling into a string type.
    jsonData := []byte("123")
    var flag feature.Flag
    err := json.Unmarshal(jsonData, &flag)
    if err == nil {
        t.Fatalf("Expected error when unmarshalling invalid JSON type into feature.Flag, got nil")
    }
}
// TestFlagNullJSONUnmarshalling tests that unmarshalling a JSON null into a feature.Flag returns an error.
func TestFlagNullJSONUnmarshalling(t *testing.T) {
    // Test that unmarshalling JSON null produces an empty feature.Flag without error.
    jsonData := []byte("null")
    var flag feature.Flag
    err := json.Unmarshal(jsonData, &flag)
    if err != nil {
        t.Fatalf("Unexpected error when unmarshalling null JSON into feature.Flag: %v", err)
    }
    if flag != "" {
        t.Errorf("Expected empty feature flag when unmarshalling JSON null, got: %q", flag)

    // Test with a flag containing Unicode characters.
    }
    unicodeFlag := feature.Flag("测试Unicode🚀")
    expected := "测试Unicode🚀"
    if string(unicodeFlag) != expected {
        t.Errorf("Expected flag %q, got %q", expected, unicodeFlag)
    }

    // Verify that JSON marshalling works correctly with Unicode.
    data, err := json.Marshal(unicodeFlag)
    if err != nil {
        t.Fatalf("Failed to marshal unicode flag: %v", err)
    }
    if string(data) != "\"测试Unicode🚀\"" {
        t.Errorf("Expected JSON output %q, got %q", "\"测试Unicode🚀\"", string(data))
    }
}
// TestFlagMapJSONMarshalling tests that a map with feature.Flag keys is correctly marshalled and then unmarshalled.
func TestFlagMapJSONMarshalling(t *testing.T) {
    // Create a map that uses feature.Flag as keys.
    flagMap := map[feature.Flag]string{
        feature.Flag("KeyOne"): "ValueOne",
        feature.Flag("KeyTwo"): "ValueTwo",
    }

    // Marshal the map to JSON. Note that JSON requires keys to be strings.
    data, err := json.Marshal(flagMap)
    if err != nil {
        t.Fatalf("Failed to marshal flag map: %v", err)
    }

    // Unmarshal the JSON data into a map[string]string.
    var unmarshalled map[string]string
    if err := json.Unmarshal(data, &unmarshalled); err != nil {
        t.Fatalf("Failed to unmarshal flag map JSON: %v", err)
    }

    // Verify that each key in the original map appears in the unmarshalled map.
    for key, originalValue := range flagMap {
        strKey := string(key)
        value, ok := unmarshalled[strKey]
        if !ok {
            t.Errorf("Key %q missing in unmarshalled map", strKey)
        }
        if value != originalValue {
            t.Errorf("For key %q, expected value %q, got %q", strKey, originalValue, value)
        }
    }
}

// TestFlagWhitespace tests that feature.Flag preserves whitespace characters.
func TestFlagWhitespace(t *testing.T) {
    // Create a flag with leading and trailing whitespace.
    whitespaceFlag := feature.Flag("   FlagWithSpaces   ")
    // Verify that converting the flag to a string preserves the whitespace.
    expected := "   FlagWithSpaces   "
    if string(whitespaceFlag) != expected {
        t.Errorf("Expected flag to be %q, got %q", expected, whitespaceFlag)
    }
}
// TestFlagFormatting tests that feature.Flag formats correctly with fmt.Sprintf.
func TestFlagFormatting(t *testing.T) {
    flag := feature.Flag("FormatTest")
    formattedS := fmt.Sprintf("%s", flag)
    formattedV := fmt.Sprintf("%v", flag)
    expected := "FormatTest"
    if formattedS != expected {
        t.Errorf("Expected fmt %%s output %q, got %q", expected, formattedS)
    }
    if formattedV != expected {
        t.Errorf("Expected fmt %%v output %q, got %q", expected, formattedV)
    }
}
// TestFlagSorting verifies that a slice of feature.Flag values can be sorted
// using their underlying string values.
func TestFlagSorting(t *testing.T) {
    flags := []feature.Flag{
        feature.Flag("zeta"),
        feature.Flag("alpha"),
        feature.Flag("gamma"),
        feature.Flag("beta"),
    }
    sort.Slice(flags, func(i, j int) bool {
        return string(flags[i]) < string(flags[j])
    })
    expected := []feature.Flag{
        feature.Flag("alpha"),
        feature.Flag("beta"),
        feature.Flag("gamma"),
        feature.Flag("zeta"),
    }
    for i, flag := range flags {
        if flag != expected[i] {
            t.Errorf("At index %d, expected flag %q, got %q", i, expected[i], flag)
        }
    }
}

// TestFlagJoin verifies that a slice of feature.Flag values can be joined together
// into a single comma-separated string.
func TestFlagJoin(t *testing.T) {
    flags := []feature.Flag{
        feature.Flag("Test1"),
        feature.Flag("Test2"),
        feature.Flag("Test3"),
    }
    // Convert each feature.Flag to its underlying string.
    strSlice := make([]string, len(flags))
    for i, flag := range flags {
        strSlice[i] = string(flag)
    }
    result := strings.Join(strSlice, ",")
    expected := "Test1,Test2,Test3"
    if result != expected {
        t.Errorf("Expected joined string %q, got %q", expected, result)
    }
}

// TestFlagTrimSpace verifies that common string functions such as TrimSpace work
// correctly on feature.Flag values.
func TestFlagTrimSpace(t *testing.T) {
    flag := feature.Flag("   trim me   ")
    // Use TrimSpace on the underlying string.
    trimmed := strings.TrimSpace(string(flag))
    expected := "trim me"
    if trimmed != expected {
        t.Errorf("Expected trimmed flag to be %q, got %q", expected, trimmed)
    }
}