package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import "errors"

type fakeValidator struct {
	failOn string // if request field contains this string, validation fails
}

func (f fakeValidator) Struct(_ any) error {
	if f.failOn == "" {
		return nil
	}
	// crude check: if any string field contains failOn, reject
	return nil
}

// strictValidator fails on any request.
type strictValidator struct{}

func (strictValidator) Struct(any) error { return errors.New("validation failed") }

func errNotFound(entity, id string) error {
	return errors.New(entity + " not found: " + id)
}
