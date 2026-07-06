package web

import "embed"

// Templates holds the HTML templates.
//
//go:embed templates/*.html
var Templates embed.FS

// Static holds client-side assets served under /js/ etc.
//
//go:embed static
var Static embed.FS
