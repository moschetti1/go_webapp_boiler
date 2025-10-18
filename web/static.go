package web

import "embed"

//go:embed static/*
var StaticFilesFS embed.FS
