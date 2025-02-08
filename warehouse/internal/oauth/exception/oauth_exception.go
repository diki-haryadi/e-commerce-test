package oauthException

import (
	errorList "github.com/diki-haryadi/go-micro-template/pkg/constant/error/error_list"
	customErrors "github.com/diki-haryadi/go-micro-template/pkg/error/custom_error"
	errorUtils "github.com/diki-haryadi/go-micro-template/pkg/error/error_utils"
)

func AuthorizationCodeGrantValidationExc(err error) error {
	ve, ie := errorUtils.ValidationErrorHandler(err)
	if ie != nil {
		return ie
	}

	validationError := errorList.InternalErrorList.ValidationError
	return customErrors.NewValidationError(validationError.Msg, validationError.Code, ve)
}

func AuthorizationCodeGrantBindingExc() error {
	oauthBindingError := errorList.InternalErrorList.OauthExceptions.BindingError
	return customErrors.NewBadRequestError(oauthBindingError.Msg, oauthBindingError.Code, nil)
}

func PasswordGrantValidationExc(err error) error {
	ve, ie := errorUtils.ValidationErrorHandler(err)
	if ie != nil {
		return ie
	}

	validationError := errorList.InternalErrorList.ValidationError
	return customErrors.NewValidationError(validationError.Msg, validationError.Code, ve)
}

func PasswordGrantBindingExc() error {
	oauthBindingError := errorList.InternalErrorList.OauthExceptions.BindingError
	return customErrors.NewBadRequestError(oauthBindingError.Msg, oauthBindingError.Code, nil)
}

func GrantClientCredentialGrantValidationExc(err error) error {
	ve, ie := errorUtils.ValidationErrorHandler(err)
	if ie != nil {
		return ie
	}

	validationError := errorList.InternalErrorList.ValidationError
	return customErrors.NewValidationError(validationError.Msg, validationError.Code, ve)
}

func GrantClientCredentialGrantBindingExc() error {
	oauthBindingError := errorList.InternalErrorList.OauthExceptions.BindingError
	return customErrors.NewBadRequestError(oauthBindingError.Msg, oauthBindingError.Code, nil)
}

func RefreshTokenGrantValidationExc(err error) error {
	ve, ie := errorUtils.ValidationErrorHandler(err)
	if ie != nil {
		return ie
	}

	validationError := errorList.InternalErrorList.ValidationError
	return customErrors.NewValidationError(validationError.Msg, validationError.Code, ve)
}

func RefreshTokenGrantBindingExc() error {
	oauthBindingError := errorList.InternalErrorList.OauthExceptions.BindingError
	return customErrors.NewBadRequestError(oauthBindingError.Msg, oauthBindingError.Code, nil)
}

func IntrospectValidationExc(err error) error {
	ve, ie := errorUtils.ValidationErrorHandler(err)
	if ie != nil {
		return ie
	}

	validationError := errorList.InternalErrorList.ValidationError
	return customErrors.NewValidationError(validationError.Msg, validationError.Code, ve)
}

func IntrospectBindingExc() error {
	oauthBindingError := errorList.InternalErrorList.OauthExceptions.BindingError
	return customErrors.NewBadRequestError(oauthBindingError.Msg, oauthBindingError.Code, nil)
}
