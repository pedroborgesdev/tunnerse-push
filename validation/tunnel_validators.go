package validation

import (
	"errors"
	"regexp"
)

var (
	ErrBannedName = errors.New("tunnel name contais inavalid characters")
)

type TunnelValidator struct {
	nameRegex *regexp.Regexp
}

func NewTunnelValidator() *TunnelValidator {
	return &TunnelValidator{
		nameRegex: regexp.MustCompile(`^[A-Za-z][A-Za-z@_-]*$`),
	}
}

func (v *TunnelValidator) ValidateTunnelRegister(name string) error {
	// if err := v.validadeName(name); err != nil {
	// 	return err
	// }

	return nil
}

// func (v *TunnelValidator) validadeName(name string) error {
// 	if !v.nameRegex.MatchString(name) {
// 		return ErrBannedName
// 	}
// 	return nil
// }
