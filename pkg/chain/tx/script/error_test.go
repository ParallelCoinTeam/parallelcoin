package txscript

import (
	"testing"
)

// TestErrorCodeStringer tests the stringized output for the ErrorCode type.
func TestErrorCodeStringer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   ErrorCode
		want string
	}{
		{errInternal, "ErrInternal"},
		{errInvalidFlags, "ErrInvalidFlags"},
		{errInvalidIndex, "ErrInvalidIndex"},
		{errUnsupportedAddress, "ErrUnsupportedAddress"},
		{errTooManyRequiredSigs, "ErrTooManyRequiredSigs"},
		{errTooMuchNullData, "ErrTooMuchNullData"},
		{errNotMultisigScript, "ErrNotMultisigScript"},
		{errEarlyReturn, "ErrEarlyReturn"},
		{errEmptyStack, "ErrEmptyStack"},
		{errEvalFalse, "ErrEvalFalse"},
		{ErrScriptUnfinished, "ErrScriptUnfinished"},
		{ErrInvalidProgramCounter, "ErrInvalidProgramCounter"},
		{ErrScriptTooBig, "ErrScriptTooBig"},
		{ErrElementTooBig, "ErrElementTooBig"},
		{ErrTooManyOperations, "ErrTooManyOperations"},
		{ErrStackOverflow, "ErrStackOverflow"},
		{ErrInvalidPubKeyCount, "ErrInvalidPubKeyCount"},
		{ErrInvalidSignatureCount, "ErrInvalidSignatureCount"},
		{errNumberTooBig, "ErrNumberTooBig"},
		{errVerify, "ErrVerify"},
		{errEqualVerify, "ErrEqualVerify"},
		{errNumEqualVerify, "ErrNumEqualVerify"},
		{errCheckSigVerify, "ErrCheckSigVerify"},
		{errCheckMultiSigVerify, "ErrCheckMultiSigVerify"},
		{errDisabledOpcode, "ErrDisabledOpcode"},
		{errReservedOpcode, "ErrReservedOpcode"},
		{errMalformedPush, "ErrMalformedPush"},
		{errInvalidStackOperation, "ErrInvalidStackOperation"},
		{errUnbalancedConditional, "ErrUnbalancedConditional"},
		{errMinimalData, "ErrMinimalData"},
		{errInvalidSigHashType, "ErrInvalidSigHashType"},
		{errSigTooShort, "ErrSigTooShort"},
		{errSigTooLong, "ErrSigTooLong"},
		{errSigInvalidSeqID, "ErrSigInvalidSeqID"},
		{ErrSigInvalidDataLen, "ErrSigInvalidDataLen"},
		{ErrSigMissingSTypeID, "ErrSigMissingSTypeID"},
		{ErrSigMissingSLen, "ErrSigMissingSLen"},
		{ErrSigInvalidSLen, "ErrSigInvalidSLen"},
		{ErrSigInvalidRIntID, "ErrSigInvalidRIntID"},
		{ErrSigZeroRLen, "ErrSigZeroRLen"},
		{ErrSigNegativeR, "ErrSigNegativeR"},
		{ErrSigTooMuchRPadding, "ErrSigTooMuchRPadding"},
		{ErrSigInvalidSIntID, "ErrSigInvalidSIntID"},
		{ErrSigZeroSLen, "ErrSigZeroSLen"},
		{ErrSigNegativeS, "ErrSigNegativeS"},
		{ErrSigTooMuchSPadding, "ErrSigTooMuchSPadding"},
		{ErrSigHighS, "ErrSigHighS"},
		{ErrNotPushOnly, "ErrNotPushOnly"},
		{ErrSigNullDummy, "ErrSigNullDummy"},
		{ErrPubKeyType, "ErrPubKeyType"},
		{ErrCleanStack, "ErrCleanStack"},
		{ErrNullFail, "ErrNullFail"},
		{ErrDiscourageUpgradableNOPs, "ErrDiscourageUpgradableNOPs"},
		{ErrNegativeLockTime, "ErrNegativeLockTime"},
		{ErrUnsatisfiedLockTime, "ErrUnsatisfiedLockTime"},
		{ErrWitnessProgramEmpty, "ErrWitnessProgramEmpty"},
		{ErrWitnessProgramMismatch, "ErrWitnessProgramMismatch"},
		{ErrWitnessProgramWrongLength, "ErrWitnessProgramWrongLength"},
		{ErrWitnessMalleated, "ErrWitnessMalleated"},
		{ErrWitnessMalleatedP2SH, "ErrWitnessMalleatedP2SH"},
		{ErrWitnessUnexpected, "ErrWitnessUnexpected"},
		{ErrMinimalIf, "ErrMinimalIf"},
		{ErrWitnessPubKeyType, "ErrWitnessPubKeyType"},
		{ErrDiscourageUpgradableWitnessProgram, "ErrDiscourageUpgradableWitnessProgram"},
		{0xffff, "Unknown ErrorCode (65535)"},
	}
	// Detect additional error codes that don't have the stringer added.
	if len(tests)-1 != int(numErrorCodes) {
		t.Errorf("It appears an error code was added without adding an " +
			"associated stringer test")
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.String()
		if result != test.want {
			t.Errorf("String #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}
}

// TestError tests the error output for the ScriptError type.
func TestError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   ScriptError
		want string
	}{
		{
			ScriptError{Description: "some error"},
			"some error",
		},
		{
			ScriptError{Description: "human-readable error"},
			"human-readable error",
		},
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.Error()
		if result != test.want {
			t.Errorf("ScriptError #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}
}
