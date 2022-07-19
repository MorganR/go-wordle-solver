module github.com/MorganR/go-wordle-solver/bin

go 1.18

replace github.com/MorganR/go-wordle-solver/lib => ../lib

require (
	github.com/MorganR/go-wordle-solver/lib v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.5.0
)

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/exp v0.0.0-20220613132600-b0d781184e0d // indirect
)
