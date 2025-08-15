module dstartest

go 1.24.5

replace github.com/derekmwright/htemel => ../htemel

require (
	github.com/derekmwright/htemel v0.0.0-20250813114536-7c3d1277f268
	github.com/go-chi/chi/v5 v5.2.2
	github.com/starfederation/datastar-go v1.0.1
)

require (
	github.com/CAFxX/httpcompression v0.0.9 // indirect
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/net v0.43.0 // indirect
)
