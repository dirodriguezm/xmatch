package validator

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValid(t *testing.T) {
	v := &Validator{
		FieldErrors:    map[string]string{},
		NonFieldErrors: []string{},
	}
	assert.True(t, v.Valid(), "Validator should be valid initially")

	v.CheckField(NotBlank("test"), "field1", "Field cannot be blank")
	assert.True(t, v.Valid(), "Validator should still be valid after checking a valid field")

	v.CheckField(NotBlank(""), "field2", "Field cannot be blank")
	assert.False(t, v.Valid(), "Validator should be invalid after checking an invalid field")

	v.AddNonFieldError("non-field error, this should make it invalid")
	assert.False(t, v.Valid(), "Validator should be invalid after adding a non-field error")
}

func TestNotBlack(t *testing.T) {
	assert.True(t, NotBlank("test"), "NotBlank should return true for non-empty string")
	assert.False(t, NotBlank(""), "NotBlank should return false for empty string")
	assert.False(t, NotBlank("   "), "NotBlank should return false for whitespace string")
}

func TestPermittedValue(t *testing.T) {
	assert.True(t, PermittedValue("a", "a", "b", "c"), "PermittedValue should return true for permitted value")
	assert.False(t, PermittedValue("d", "a", "b", "c"), "PermittedValue should return false for non-permitted value")
}

func TestMinChars(t *testing.T) {
	assert.True(t, MinChars("test", 4), "MinChars should return true for string with enough characters")
	assert.False(t, MinChars("test", 5), "MinChars should return false for string with not enough characters")
}

func TestMaxChars(t *testing.T) {
	assert.True(t, MaxChars("test", 4), "MaxChars should return true for string with enough characters")
	assert.False(t, MaxChars("test", 3), "MaxChars should return false for string with too many characters")
}

func TestCheckField(t *testing.T) {
	v := &Validator{}
	v.CheckField(NotBlank("test"), "field1", "Field cannot be blank")
	assert.Empty(t, v.FieldErrors, "FieldErrors should be empty for valid field")

	v.CheckField(NotBlank(""), "field2", "Field cannot be blank")
	assert.NotEmpty(t, v.FieldErrors, "FieldErrors should contain an error for invalid field")
	assert.Equal(t, "Field cannot be blank", v.FieldErrors["field2"], "Error message should match")
}

func CheckAddFieldError(t *testing.T) {
	v := &Validator{}
	v.AddFieldError("field1", "Error message for field1")
	assert.Equal(t, "Error message for field1", v.FieldErrors["field1"], "FieldErrors should contain the added error")

	// Adding the same error again should not change the existing error
	v.AddFieldError("field1", "Another error message for field1")
	assert.Equal(t, "Error message for field1", v.FieldErrors["field1"], "FieldErrors should not change on duplicate key")
}

func CheckAddNonFieldError(t *testing.T) {
	v := &Validator{}
	v.AddNonFieldError("Non-field error message")
	assert.Equal(t, 1, len(v.NonFieldErrors), "NonFieldErrors should contain one error")
	assert.Equal(t, "Non-field error message", v.NonFieldErrors[0], "NonFieldErrors should contain the added error")

	// Adding another non-field error
	v.AddNonFieldError("Another non-field error message")
	assert.Equal(t, 2, len(v.NonFieldErrors), "NonFieldErrors should contain two errors")
	assert.Equal(t, "Another non-field error message", v.NonFieldErrors[1], "NonFieldErrors should contain the second added error")
}

func TestMatches(t *testing.T) {
	rx := regexp.MustCompile("^[a-z]+$")

	assert.True(t, Matches("test", rx), "Matches should return true for string matching regex")
	assert.False(t, Matches("Test123", rx), "Matches should return false for string not matching regex")
	assert.False(t, Matches("", rx), "Matches should return false for empty string")
}

func TestEmailRX(t *testing.T) {
	assert.True(t, Matches("mati@gmail.com", EmailRX), "EmailRX should match valid email")
	assert.False(t, Matches("invalid-email", EmailRX), "EmailRX should not match invalid email")
	assert.False(t, Matches("test@.com", EmailRX), "EmailRX should not match email with missing domain")
	assert.False(t, Matches("test@domain", EmailRX), "EmailRX should not match email with missing TLD")
	assert.False(t, Matches("test@domain.", EmailRX), "EmailRX should not match email with trailing dot in domain")
	assert.False(t, Matches("test@domain..com", EmailRX), "EmailRX should not match email with double dot in domain")
}
