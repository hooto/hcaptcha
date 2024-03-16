// Copyright 2015 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
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
	"github.com/hooto/httpsrv"
)

func WebServerModule() *httpsrv.Module {

	mod := httpsrv.NewModule()

	mod.RegisterController(new(Api))

	return mod
}

func WebServerStart() {

	httpsrv.DefaultService.Config.HttpPort = gcfg.ServerPort

	httpsrv.DefaultService.HandleModule("/hcaptcha", WebServerModule())

	httpsrv.DefaultService.Start()
}

type Api struct {
	*httpsrv.Controller
}

func (c Api) VerifyAction() {

	if err := Verify(c.Params.Value("hcaptcha_token"),
		c.Params.Value("hcaptcha_word")); err != nil {
		c.RenderString("false\n" + err.Code)
	} else {
		c.RenderString("true")
	}
}

func (c Api) ImageAction() {

	c.AutoRender = false

	reload := false

	if c.Params.Value("hcaptcha_opt") == "refresh" {
		reload = true
	}

	if img, err := ImageFetch(c.Params.Value("hcaptcha_token"), reload); err != nil {
		c.RenderError(500, err.Code)
	} else {
		c.Response.Out.Header().Set("Content-type", "image/png")
		c.Response.Out.Write(img)
	}
}
