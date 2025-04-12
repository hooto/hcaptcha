// Copyright 2015 Eryx <evorui at gmail dot com>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package captcha4g

import (
	"image"
)

type FontEntry struct {
	Symbol string
	Width  int
	Height int
	Image  *image.RGBA
}

type FontList struct {
	MaxHeight int
	Length    int
	Items     []*FontEntry
}

// ErrorMeta provides more information about an api failure.
type ErrorMeta struct {
	// A machine-readable description of the type of the error. If this value is
	// empty there is no information available.
	Code string `json:"code,omitempty" toml:"code,omitempty"`

	// A human-readable description of the error message.
	Message string `json:"message,omitempty" toml:"message,omitempty"`
}

const (
	ErrCodeBadArgument  = "BadArgument"
	ErrCodeUnavailable  = "Unavailable"
	ErrCodeServerError  = "ServerError"
	ErrCodeNotFound     = "NotFound"
	ErrCodeAccessDenied = "AccessDenied"
	ErrCodeUnauthorized = "Unauthorized"
)

func NewErrorMeta(code, message string) *ErrorMeta {
	return &ErrorMeta{
		Code:    code,
		Message: message,
	}
}
