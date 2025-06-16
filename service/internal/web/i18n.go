// Copyright 2024-2025 Matías Medina Silva
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package web

import (
	"fmt"
	"io/fs"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/dirodriguezm/xmatch/service/ui"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	defaultLang = language.English
)

var translations *i18n.Bundle
var defaultLocalizer *i18n.Localizer

func loadTranslations() error {
	bundle := i18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	locales, err := fs.Glob(ui.Files, "locale/*.toml")
	if err != nil {
		return fmt.Errorf("could not load locales: %v", err)
	}
	for _, locale := range locales {
		_, err = bundle.LoadMessageFileFS(ui.Files, locale)
		if err != nil {
			return fmt.Errorf("error loading bundle locale file %s: %v", locale, err)
		}
	}

	translations = bundle
	defaultLocalizer = i18n.NewLocalizer(bundle, defaultLang.String())
	return nil
}

func translate(localizer *i18n.Localizer, id string, args ...any) string {
	var data map[string]any
	if len(args) > 0 {
		data = make(map[string]any, len(args))
		for n, iface := range args {
			data["v"+strconv.Itoa(n)] = iface
		}
	}
	str, _, err := localizer.LocalizeWithTag(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: data,
	})
	if str == "" && err != nil {
		return "[TL err: " + err.Error() + "]"
	}
	return str
}

func translateCount(localizer *i18n.Localizer, id string, ct int, args ...any) string {
	data := make(map[string]any, len(args)+1)
	if len(args) > 0 {
		for n, iface := range args {
			data["v"+strconv.Itoa(n)] = iface
		}
	}
	data["ct"] = ct
	str, _, err := localizer.LocalizeWithTag(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: data,
		PluralCount:  ct,
	})
	if str == "" && err != nil {
		return "[TL err: " + err.Error() + "]"
	}
	return str
}
