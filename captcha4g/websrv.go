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

// This file is the github.com/hooto/httpsrv/v2 entry point for hcaptcha. The
// optional gofiber/fiber/v3 entry point lives in captcha4g/webfiber. Both
// coexist so hcaptcha can be served by either framework.

package captcha4g

import (
	"fmt"

	"github.com/hooto/httpsrv/v2"
)

// Register mounts the hcaptcha API routes on an httpsrv v2 router, reproducing
// the routes previously exposed by the v1 "Api" controller:
//
//	{prefix}/api/verify, {prefix}/api/image
//
// The caller mounts the router at the module prefix (e.g. "/hcaptcha"). Routes
// accept any HTTP method, matching the method-agnostic v1 dispatch.
func Register(router httpsrv.Router) {
	router.All("/api/verify", srvApiVerify)
	router.All("/api/image", srvApiImage)
}

// WebServerStart starts a standalone httpsrv v2 server bound to the configured
// port and mounts the hcaptcha API under "/hcaptcha". It blocks until the server
// stops.
func WebServerStart() {
	app := httpsrv.New(httpsrv.WithConfig(httpsrv.Config{
		Addr: fmt.Sprintf(":%d", gcfg.ServerPort),
	}))
	Register(app.Group("/hcaptcha"))
	_ = app.Run()
}

// srvParam replicates the v1 Params.Value precedence (path -> query -> form)
// for the captcha request parameters, which arrive as query/form values.
func srvParam(c httpsrv.Ctx, key string) string {
	if v := c.Params(key); v != "" {
		return v
	}
	if v := c.Query(key); v != "" {
		return v
	}
	return c.FormValue(key)
}

func srvApiVerify(c httpsrv.Ctx) error {
	if e := Verify(srvParam(c, "hcaptcha_token"), srvParam(c, "hcaptcha_word")); e != nil {
		return c.SendString("false\n" + e.Code)
	}
	return c.SendString("true")
}

func srvApiImage(c httpsrv.Ctx) error {
	reload := srvParam(c, "hcaptcha_opt") == "refresh"

	img, e := ImageFetch(srvParam(c, "hcaptcha_token"), reload)
	if e != nil {
		c.Status(500)
		return c.SendString(e.Code)
	}
	c.SetHeader("Content-type", "image/png")
	return c.Send(img)
}
