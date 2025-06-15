package web

import (
	"fmt"
	"io/fs"
	"log/slog"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/dirodriguezm/xmatch/service/ui"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	defaultLang = language.Spanish
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
		slog.Debug("locale loaded", "locale", locale)
	}

	translations = bundle
	defaultLocalizer = i18n.NewLocalizer(bundle, defaultLang.String())
	return nil
}

func translateFunc(localizer *i18n.Localizer, id string, args ...any) string {
	var data map[string]any
	if len(args) > 0 {
		data = make(map[string]any, len(args))
		for n, iface := range args {
			data["v"+strconv.Itoa(n)] = iface
		}
	}
	slog.Debug("translation", "id", id, "args", args)
	str, _, err := localizer.LocalizeWithTag(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: data,
	})
	if str == "" && err != nil {
		return "[TL err: " + err.Error() + "]"
	}
	return str
}

func translateCountFunc(localizer *i18n.Localizer) func(string, int, ...any) string {
	return func(id string, ct int, args ...any) string {
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
}
