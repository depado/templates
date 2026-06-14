module {{ .gitserver }}/{{ .owner }}/{{ .name }}

go 1.26

require (
	github.com/pocketbase/pocketbase v0.39.4
	github.com/spf13/viper v1.21.0
)
