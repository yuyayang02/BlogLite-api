package errors_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yuyayang02/BlogLite-api/internal/common/errors"
	"testing"
)

func TestError_Wrap(t *testing.T) {
	assert.NoError(t, errors.Wrap(nil, errors.ErrorCodeDomainError, "test nil"))

	appError := errors.New(errors.ErrorCodeDomainError, "test")
	err := errors.Wrap(appError, errors.ErrorCodeDomainError, "test2")

	assert.Error(t, err)
}

func TestError_Cause(t *testing.T) {
	assert.NoError(t, errors.Cause(nil))

	appError := fmt.Errorf("test")
	err2 := errors.Wrap(appError, errors.ErrorCodeResourceNotFound, "test2")
	err := errors.Wrap(err2, errors.ErrorCodeResourceNotFound, "test3")

	assert.Equal(t, errors.Cause(err), appError)

	appError = errors.New(errors.ErrorCodeDomainError, "test")
	err2 = errors.Wrap(appError, errors.ErrorCodeResourceNotFound, "test2")
	err = errors.Wrap(err2, errors.ErrorCodeResourceNotFound, "test3")

	assert.Equal(t, errors.Cause(err), appError)
}

func TestError_Reason(t *testing.T) {
	assert.NoError(t, errors.Reason(nil))

	appError := fmt.Errorf("test")
	err2 := errors.Wrap(appError, errors.ErrorCodeResourceNotFound, "test2")
	err := errors.Wrap(err2, errors.ErrorCodeResourceNotFound, "test3")

	assert.Equal(t, errors.Reason(err), err2)

	appError = errors.New(errors.ErrorCodeDomainError, "test")
	err2 = errors.Wrap(appError, errors.ErrorCodeResourceNotFound, "test2")
	err = errors.Wrap(err2, errors.ErrorCodeResourceNotFound, "test3")

	assert.Equal(t, errors.Reason(err), appError)
}
