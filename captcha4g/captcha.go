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
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"math"
	"math/rand"
)

func Verify(token, word string) *ErrorMeta {

	if token == "" || word == "" {
		return &ErrorMeta{"invalid-request", ""}
	}

	if DataConnector == nil {
		return &ErrorMeta{"hcaptcha-not-reachable", ""}
	}

	//
	if rs := DataConnector.NewReader(_token_word_key(token)).Exec(); !rs.OK() ||
		rs.Item().StringValue() != word {
		return &ErrorMeta{"incorrect-hcaptcha-word", ""}
	}

	DataConnector.NewWriter(_token_word_key(token), _token_image_key(token)).Exec()

	return nil
}

func ImageFetch(token string, reload bool) ([]byte, *ErrorMeta) {

	if DataConnector == nil {
		return []byte{}, &ErrorMeta{"hcaptcha-not-reachable", ""}
	}

	if !reload {
		if rs := DataConnector.NewReader(_token_image_key(token)).Exec(); rs.OK() {
			return rs.Item().Value, nil
		}
	}

	vylen := gcfg.LengthMin + rand.Intn(gcfg.LengthMax-gcfg.LengthMin+1)

	capstr := image.NewRGBA(image.Rect(0, 0, gcfg.ImageWidth, gcfg.ImageHeight))

	prev_min_x, prev_min_y, prev_max_x, prev_max_y := 0, 0, 0, 0

	vyword := ""

	for i := 0; i < vylen; i++ {

		font := fonts.Items[rand.Intn(fonts.Length)]

		yshift := rand.Intn(int(float64(font.Height) * (1 - (2 * gcfg.fluctuation_amplitude))))

		start := gcfg.font_size - int(float64(font.Height)*(1-gcfg.fluctuation_amplitude))

		var r image.Rectangle

		if i == 0 {

			prev_min_x, prev_min_y = gcfg.font_padding, start+yshift
			prev_max_x, prev_max_y = prev_min_x+font.Width, prev_min_y+font.Height

			r = image.Rect(prev_min_x, prev_min_y, prev_max_x, prev_max_y)

		} else {

			x, y := prev_max_x, start+yshift

			for sx := 1; sx < font.Width; sx += 1 {

				for sy := 1; sy < font.Height; sy += 1 {

					if _, _, _, a := font.Image.At(sx, sy).RGBA(); a < 5 {
						continue
					}

					target_x, target_y := prev_max_x-sx, start+yshift+sy

					if _, _, _, al := capstr.At(target_x, target_y).RGBA(); al < 5 {
						continue
					}

					x = target_x

					break
				}

				if x != prev_max_x {
					break
				}
			}

			prev_max_x = x + font.Width

			r = image.Rect(x, y, prev_max_x, y+font.Height)
		}

		if prev_max_x > (gcfg.ImageWidth - 10) {
			break
		}

		vyword += font.Symbol
		draw.Draw(capstr, r, font.Image, image.Pt(0, 0), draw.Over)
	}

	capwave := image.NewRGBA(image.Rect(0, 0, gcfg.ImageWidth, gcfg.ImageHeight))

	amplude := _rand_float(5, 10)
	period := _rand_float(100, 200)

	dx := 2.5 * math.Pi / period

	for x := 0; x < gcfg.ImageWidth; x++ {

		for y := 0; y < gcfg.ImageHeight; y++ {

			sx := x + int(amplude*math.Sin(float64(y)*dx))
			sy := y + int(amplude*math.Cos(float64(x)*dx))

			if sx < 0 || sy < 0 || sx >= gcfg.ImageWidth-1 || sy >= gcfg.ImageHeight-1 {
				continue
			}

			if capstr.RGBAAt(sx, sy).A < 1 {
				continue
			}

			capwave.Set(x, y, capstr.At(sx, sy))
		}
	}

	buf := new(bytes.Buffer)

	if err := png.Encode(buf, capwave); err != nil {
		return []byte{}, &ErrorMeta{"ServerError", err.Error()}
	}

	if rs := DataConnector.NewWriter(_token_word_key(token), []byte(vyword)).
		SetTTL(gcfg.ImageExpiration).Exec(); !rs.OK() {
		return []byte{}, &ErrorMeta{"ServerError " + rs.ErrorMessage(), ""}
	}

	if rs := DataConnector.NewWriter(_token_image_key(token), buf.Bytes()).
		SetTTL(gcfg.ImageExpiration).Exec(); !rs.OK() {
		return []byte{}, &ErrorMeta{"ServerError " + rs.ErrorMessage(), ""}
	}

	return buf.Bytes(), nil
}
