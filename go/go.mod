module {{ .gitserver }}/{{ .owner }}/{{ .name }}

go 1.22.2

require (
	{{ if .gin -}}
	github.com/Depado/ginprom v1.8.3
	github.com/gin-contrib/cors v1.7.7
	github.com/gin-gonic/gin v1.12.0
	{{ end -}}
	github.com/rs/zerolog v1.35.0
	github.com/spf13/cobra v1.10.2
	github.com/spf13/viper v1.21.0
)
