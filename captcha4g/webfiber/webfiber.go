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

// This file adds an OPTIONAL gofiber/fiber/v3 entry point (Register) alongside
// the existing httpsrv WebServerModule. Both coexist so hcaptcha can be served
// by either framework. The httpsrv support in websrv.go is unchanged.

package webfiber

import (
	"github.com/gofiber/fiber/v3"

	"github.com/hooto/hcaptcha/captcha4g"
)

// Register mounts the hcaptcha API routes on a fiber router, reproducing the
// httpsrv "Api" controller routes:
//
//	{prefix}/api/verify, {prefix}/api/image
//
// The caller mounts the router at the module prefix (e.g. "/hp/+/hcaptcha").
// Routes accept any HTTP method (httpsrv dispatch is method-agnostic).
func Register(router fiber.Router) {
	router.All("/api/verify", fiberApiVerify)
	router.All("/api/image", fiberApiImage)
}

// fiberParam replicates httpsrv Params.Value precedence (path → query → form)
// for the captcha request parameters, which arrive as query/form values.
func fiberParam(c fiber.Ctx, key string) string {
	if v := c.Params(key); v != "" {
		return v
	}
	if v := c.Query(key); v != "" {
		return v
	}
	return c.FormValue(key)
}

func fiberApiVerify(c fiber.Ctx) error {
	if err := captcha4g.Verify(fiberParam(c, "hcaptcha_token"), fiberParam(c, "hcaptcha_word")); err != nil {
		return c.SendString("false\n" + err.Code)
	}
	return c.SendString("true")
}

func fiberApiImage(c fiber.Ctx) error {
	reload := fiberParam(c, "hcaptcha_opt") == "refresh"

	img, err := captcha4g.ImageFetch(fiberParam(c, "hcaptcha_token"), reload)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		_, e := c.Write([]byte(err.Code))
		return e
	}
	c.Set("Content-type", "image/png")
	_, e := c.Write(img)
	return e
}
