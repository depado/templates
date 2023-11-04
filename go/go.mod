module {{ .gitserver }}/{{ .owner }}/{{ .name }}

go 1.21.1

require (
	{{ if .gin -}}
	github.com/Depado/ginprom v1.7.11
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-gonic/gin v1.9.1
	{{ end -}}
	github.com/rs/zerolog v1.31.0
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.17.0
)
