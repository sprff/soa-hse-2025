module social/apiservice

go 1.24.0

require (
	github.com/go-chi/chi/v5 v5.2.1
	github.com/go-chi/render v1.0.3
	social/shared v0.0.0
)

require github.com/ajg/form v1.5.1 // indirect

replace social/shared => ../../shared/source
