package errors

import (
	"testing"

	"github.com/radekg/yugabyte-db-go-client/utils"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
	"github.com/stretchr/testify/assert"
)

func TestCDCErrors(t *testing.T) {

	t.Run("it=handles nil errors", func(tt *testing.T) {
		assert.Nil(tt, NewCDCError(nil))
	})

	t.Run("it=handles errors with nil code", func(tt *testing.T) {
		anError := NewCDCError(&ybApi.CDCErrorPB{})
		typedError, ok := anError.(*CDCError)
		assert.True(tt, ok, "expected the error to be *CDCError")
		expectedErrorString := "CDC rpc error: code: 1 (UNKNOWN_ERROR), status: <unknown>"
		assert.Equal(tt, expectedErrorString, typedError.Error())
	})

	t.Run("it=handles errors with nil status", func(tt *testing.T) {
		anError := NewCDCError(&ybApi.CDCErrorPB{
			Code: utils.PCDCErrorCode(ybApi.CDCErrorPB_INVALID_REQUEST),
		})
		typedError, ok := anError.(*CDCError)
		assert.True(tt, ok, "expected the error to be *CDCError")
		expectedErrorString := "CDC rpc error: code: 10 (INVALID_REQUEST), status: <unknown>"
		assert.Equal(tt, expectedErrorString, typedError.Error())
	})

	t.Run("it=handles errors with all details", func(tt *testing.T) {
		anError := NewCDCError(&ybApi.CDCErrorPB{
			Code: utils.PCDCErrorCode(ybApi.CDCErrorPB_INVALID_REQUEST),
			Status: &ybApi.AppStatusPB{
				Code:       utils.PAppStatusErrorCode(ybApi.AppStatusPB_ABORTED),
				Message:    utils.PString("test error"),
				SourceFile: utils.PString("errors_test.go"),
				SourceLine: utils.PInt32(42),
			},
		})
		typedError, ok := anError.(*CDCError)
		assert.True(tt, ok, "expected the error to be *CDCError")
		expectedErrorString := "CDC rpc error: code: 10 (INVALID_REQUEST), status: 11 (ABORTED)\n\tmessage: test error\n\tsource: errors_test.go@42"
		assert.Equal(tt, expectedErrorString, typedError.Error())
	})

	t.Run("it=handles errors without message", func(tt *testing.T) {
		anError := NewCDCError(&ybApi.CDCErrorPB{
			Code: utils.PCDCErrorCode(ybApi.CDCErrorPB_INVALID_REQUEST),
			Status: &ybApi.AppStatusPB{
				Code:       utils.PAppStatusErrorCode(ybApi.AppStatusPB_ABORTED),
				SourceFile: utils.PString("errors_test.go"),
				SourceLine: utils.PInt32(42),
			},
		})
		typedError, ok := anError.(*CDCError)
		assert.True(tt, ok, "expected the error to be *CDCError")
		expectedErrorString := "CDC rpc error: code: 10 (INVALID_REQUEST), status: 11 (ABORTED)\n\tsource: errors_test.go@42"
		assert.Equal(tt, expectedErrorString, typedError.Error())
	})

	t.Run("it=handles errors without source line", func(tt *testing.T) {
		anError := NewCDCError(&ybApi.CDCErrorPB{
			Code: utils.PCDCErrorCode(ybApi.CDCErrorPB_INVALID_REQUEST),
			Status: &ybApi.AppStatusPB{
				Code:       utils.PAppStatusErrorCode(ybApi.AppStatusPB_ABORTED),
				Message:    utils.PString("test error"),
				SourceFile: utils.PString("errors_test.go"),
			},
		})
		typedError, ok := anError.(*CDCError)
		assert.True(tt, ok, "expected the error to be *CDCError")
		expectedErrorString := "CDC rpc error: code: 10 (INVALID_REQUEST), status: 11 (ABORTED)\n\tmessage: test error\n\tsource: errors_test.go"
		assert.Equal(tt, expectedErrorString, typedError.Error())
	})

	t.Run("it=handles errors without source information", func(tt *testing.T) {
		anError := NewCDCError(&ybApi.CDCErrorPB{
			Code: utils.PCDCErrorCode(ybApi.CDCErrorPB_INVALID_REQUEST),
			Status: &ybApi.AppStatusPB{
				Code:       utils.PAppStatusErrorCode(ybApi.AppStatusPB_ABORTED),
				Message:    utils.PString("test error"),
				SourceLine: utils.PInt32(42),
			},
		})
		typedError, ok := anError.(*CDCError)
		assert.True(tt, ok, "expected the error to be *CDCError")
		expectedErrorString := "CDC rpc error: code: 10 (INVALID_REQUEST), status: 11 (ABORTED)\n\tmessage: test error"
		assert.Equal(tt, expectedErrorString, typedError.Error())
	})

}
