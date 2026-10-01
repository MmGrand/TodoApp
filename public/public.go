// Package public встраивает статические файлы фронтенда в бинарник,
// чтобы приложение не зависело от рабочей директории и PROJECT_ROOT.
package public

import "embed"

//go:embed index.html
var FS embed.FS
