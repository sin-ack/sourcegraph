package languages

import (
	"github.com/go-enry/go-enry/v2"
	"golang.org/x/exp/slices"
)

// TODO: We probably want to move the config for language detection into here from the syntax highlighting part
// I didn't add that yet.

// GetLanguage returns the language for the given path and contents.
func GetLanguage(path, contents string) (lang string, found bool) {
	// Force the use of the shebang.
	if shebangLang, ok := overrideViaShebang(path, contents); ok {
		return shebangLang, true
	}

	// Lastly, fall back to whatever enry decides is a useful algorithm for calculating.
	lang = enry.GetLanguage(path, []byte(contents))
	if lang != "" {
		return lang, true
	}

	return lang, false
}

// overrideViaShebang handles explicitly using the shebang whenever possible.
//
// It also covers some edge cases when enry eagerly returns more languages
// than necessary, which ends up overriding the shebang completely (which,
// IMO is the highest priority match we can have).
//
// For example, enry will return "Perl" and "Pod" for a shebang of `#!/usr/bin/env perl`.
// This is actually unhelpful, because then enry will *not* select "Perl" as the
// language (which is our desired behavior).
func overrideViaShebang(path, content string) (lang string, ok bool) {
	shebangs := enry.GetLanguagesByShebang(path, []byte(content), []string{})
	if len(shebangs) == 0 {
		return "", false
	}

	if len(shebangs) == 1 {
		return shebangs[0], true
	}

	// There are some shebangs that enry returns that are not really
	// useful for our syntax highlighters to distinguish between.
	if slices.Equal(shebangs, []string{"Perl", "Pod"}) {
		return "Perl", true
	}

	return "", false
}
