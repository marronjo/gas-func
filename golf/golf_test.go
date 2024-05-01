package golf

import (
	"testing"
)

func TestMissingFirstBracket(t *testing.T) {
	missingBracket := "deposit)"
	result, err := SearchFuncSelector(missingBracket)

	if result != (Result{}) || err == nil {
		t.Errorf("Expected error from missing bracket input, none found")
	}
	log(t, err, missingBracket)
}

func TestMissingSecondBracket(t *testing.T) {
	missingBracket := "mint("
	result, err := SearchFuncSelector(missingBracket)

	if result != (Result{}) || err == nil {
		t.Errorf("Expected error from missing bracket input, none found")
	}
	log(t, err, missingBracket)
}

func TestInvalidFuncName(t *testing.T) {
	invalidName := "transf£er()"
	result, err := SearchFuncSelector(invalidName)

	if result != (Result{}) || err == nil {
		t.Errorf("Expected error from invalid name input, none found")
	}
	log(t, err, invalidName)
}

func TestTooManyBracketsStart(t *testing.T) {
	invalidName := "burn(()"
	result, err := SearchFuncSelector(invalidName)

	if result != (Result{}) || err == nil {
		t.Errorf("Expected error from invalid number of brackets, none found")
	}
	log(t, err, invalidName)
}

func TestTooManyBracketsEnd(t *testing.T) {
	invalidName := "transfer())"
	result, err := SearchFuncSelector(invalidName)

	if result != (Result{}) || err == nil {
		t.Errorf("Expected error from invalid number of brackets, none found")
	}
	log(t, err, invalidName)
}

func TestInvalidBracketOrder(t *testing.T) {
	invalidName := "transfer)("
	result, err := SearchFuncSelector(invalidName)

	if result != (Result{}) || err == nil {
		t.Errorf("Expected error from invalid bracket order, none found")
	}
	log(t, err, invalidName)
}

func TestInvalidUint(t *testing.T) {
	invalidName := "transfer(uint7)"
	result, err := SearchFuncSelector(invalidName)

	if result != (Result{}) || err == nil {
		t.Errorf("Expected error from invalid input type, none found")
	}
	log(t, err, invalidName)
}

func TestInvalidInt(t *testing.T) {
	invalidName := "transfer(int254a)"
	result, err := SearchFuncSelector(invalidName)

	if result != (Result{}) || err == nil {
		t.Errorf("Expected error from invalid input type, none found")
	}
	log(t, err, invalidName)
}

func TestInvalidString(t *testing.T) {
	invalidName := "deposit(string1)"
	result, err := SearchFuncSelector(invalidName)

	if result != (Result{}) || err == nil {
		t.Errorf("Expected error from invalid input type, none found")
	}
	log(t, err, invalidName)
}

func log(t *testing.T, err error, input string) {
	t.Logf("Got error '%s' from input '%s'", err, input)
}
