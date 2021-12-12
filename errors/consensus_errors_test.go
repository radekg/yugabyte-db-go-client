package errors

import (
	"testing"

	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
	"github.com/stretchr/testify/assert"
)

func TestConsensusErrors(t *testing.T) {

	t.Run("it=handles nil errors", func(tt *testing.T) {
		assert.Nil(tt, NewConsensusError(nil))
	})

	t.Run("it=handles errors with nil code", func(tt *testing.T) {
		anError := NewConsensusError(&ybApi.ConsensusErrorPB{})
		typedError, ok := anError.(*ConsensusError)
		assert.True(tt, ok, "expected the error to be *ConsensusError")
		expectedErrorString := "consensus rpc error: code: 0 (UNKNOWN), status: <unknown>"
		assert.Equal(tt, expectedErrorString, typedError.Error())
	})

	t.Run("it=handles errors with nil status", func(tt *testing.T) {
		anError := NewConsensusError(&ybApi.ConsensusErrorPB{
			Code: utils.PConsensusErrorCode(ybApi.ConsensusErrorPB_INVALID_TERM),
		})
		typedError, ok := anError.(*ConsensusError)
		assert.True(tt, ok, "expected the error to be *ConsensusError")
		expectedErrorString := "consensus rpc error: code: 2 (INVALID_TERM), status: <unknown>"
		assert.Equal(tt, expectedErrorString, typedError.Error())
	})

	t.Run("it=handles errors with all details", func(tt *testing.T) {
		anError := NewConsensusError(&ybApi.ConsensusErrorPB{
			Code: utils.PConsensusErrorCode(ybApi.ConsensusErrorPB_INVALID_TERM),
			Status: &ybApi.AppStatusPB{
				Code:       utils.PAppStatusErrorCode(ybApi.AppStatusPB_ABORTED),
				Message:    utils.PString("test error"),
				SourceFile: utils.PString("errors_test.go"),
				SourceLine: utils.PInt32(42),
			},
		})
		typedError, ok := anError.(*ConsensusError)
		assert.True(tt, ok, "expected the error to be *ConsensusError")
		expectedErrorString := "consensus rpc error: code: 2 (INVALID_TERM), status: 11 (ABORTED)\n\tmessage: test error\n\tsource: errors_test.go@42"
		assert.Equal(tt, expectedErrorString, typedError.Error())
	})

	t.Run("it=handles errors without message", func(tt *testing.T) {
		anError := NewConsensusError(&ybApi.ConsensusErrorPB{
			Code: utils.PConsensusErrorCode(ybApi.ConsensusErrorPB_INVALID_TERM),
			Status: &ybApi.AppStatusPB{
				Code:       utils.PAppStatusErrorCode(ybApi.AppStatusPB_ABORTED),
				SourceFile: utils.PString("errors_test.go"),
				SourceLine: utils.PInt32(42),
			},
		})
		typedError, ok := anError.(*ConsensusError)
		assert.True(tt, ok, "expected the error to be *ConsensusError")
		expectedErrorString := "consensus rpc error: code: 2 (INVALID_TERM), status: 11 (ABORTED)\n\tsource: errors_test.go@42"
		assert.Equal(tt, expectedErrorString, typedError.Error())
	})

	t.Run("it=handles errors without source line", func(tt *testing.T) {
		anError := NewConsensusError(&ybApi.ConsensusErrorPB{
			Code: utils.PConsensusErrorCode(ybApi.ConsensusErrorPB_INVALID_TERM),
			Status: &ybApi.AppStatusPB{
				Code:       utils.PAppStatusErrorCode(ybApi.AppStatusPB_ABORTED),
				Message:    utils.PString("test error"),
				SourceFile: utils.PString("errors_test.go"),
			},
		})
		typedError, ok := anError.(*ConsensusError)
		assert.True(tt, ok, "expected the error to be *ConsensusError")
		expectedErrorString := "consensus rpc error: code: 2 (INVALID_TERM), status: 11 (ABORTED)\n\tmessage: test error\n\tsource: errors_test.go"
		assert.Equal(tt, expectedErrorString, typedError.Error())
	})

	t.Run("it=handles errors without source information", func(tt *testing.T) {
		anError := NewConsensusError(&ybApi.ConsensusErrorPB{
			Code: utils.PConsensusErrorCode(ybApi.ConsensusErrorPB_INVALID_TERM),
			Status: &ybApi.AppStatusPB{
				Code:       utils.PAppStatusErrorCode(ybApi.AppStatusPB_ABORTED),
				Message:    utils.PString("test error"),
				SourceLine: utils.PInt32(42),
			},
		})
		typedError, ok := anError.(*ConsensusError)
		assert.True(tt, ok, "expected the error to be *ConsensusError")
		expectedErrorString := "consensus rpc error: code: 2 (INVALID_TERM), status: 11 (ABORTED)\n\tmessage: test error"
		assert.Equal(tt, expectedErrorString, typedError.Error())
	})

}
