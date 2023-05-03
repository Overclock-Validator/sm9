//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"go/format"
	"io"
	"log"
	"os"
	"os/exec"
)

// Running this generator requires addchain v0.4.0, which can be installed with
//
//   go install github.com/mmcloughlin/addchain/cmd/addchain@v0.4.0
//

func generate(template, exp, element string) ([]byte, error) {
	tmplAddchainFileInvert, err := os.CreateTemp("", "addchain-template")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmplAddchainFileInvert.Name())
	if _, err := io.WriteString(tmplAddchainFileInvert, template); err != nil {
		return nil, err
	}
	if err := tmplAddchainFileInvert.Close(); err != nil {
		return nil, err
	}

	f, err := os.CreateTemp("", "addchain-gfp")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())
	cmd := exec.Command("addchain", "search", exp)
	cmd.Stderr = os.Stderr
	cmd.Stdout = f
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}
	cmd = exec.Command("addchain", "gen", "-tmpl", tmplAddchainFileInvert.Name(), f.Name())
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	out = bytes.Replace(out, []byte("Element"), []byte(element), -1)
	return format.Source(out)
}

func writeFile(fileName string, buffers ...[]byte) error {
	log.Printf("Generating %v...", fileName)
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for _, buffer := range buffers {
		if _, err := f.Write(buffer); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	out, err := generate(tmplAddchainInvert, "0xb640000002a3a6f1d603ab4ff58ec74521f2934b1a7aeedbe56f9b27e351457b", "gfP")
	if err != nil {
		log.Fatal(err)
	}
	out1, err := generate(tmplAddchainSqrt, "0x16c80000005474de3ac07569feb1d8e8a43e5269634f5ddb7cadf364fc6a28af", "gfP")
	if err != nil {
		log.Fatal(err)
	}
	if err = writeFile("gfp_invert_sqrt.go", out, out1); err != nil {
		log.Fatal(err)
	}

	out, err = generate(tmplAddchainExp1, "0x2d90000000a8e9bc7580ead3fd63b1d1487ca4d2c69ebbb6f95be6c9f8d4515f", "gfP2")
	if err != nil {
		log.Fatal(err)
	}
	out1, err = generate(tmplAddchainExp2, "0xb640000002a3a6f1d603ab4ff58ec74521f2934b1a7aeedbe56f9b27e351457d", "gfP2")
	if err != nil {
		log.Fatal(err)
	}
	out2, err := generate(tmplAddchainExp3, "0x5b2000000151d378eb01d5a7fac763a290f949a58d3d776df2b7cd93f1a8a2be", "gfP2")
	if err != nil {
		log.Fatal(err)
	}
	if err = writeFile("gfp2_sqrt.go", out, out1, out2); err != nil {
		log.Fatal(err)
	}
}

const tmplAddchainExp1 = `// Code generated by {{ .Meta.Name }}. DO NOT EDIT.
package bn256

func (e *Element) expPMinus1Over4(x *Element) *Element {
	// The sequence of {{ .Ops.Adds }} multiplications and {{ .Ops.Doubles }} squarings is derived from the
	// following addition chain generated with {{ .Meta.Module }} {{ .Meta.ReleaseTag }}.
	//
	{{- range lines (format .Script) }}
	//	{{ . }}
	{{- end }}
	//
	var z = new(Element).Set(e)
	{{- range .Program.Temporaries }}
	var {{ . }} = new(Element)
	{{- end }}
	{{ range $i := .Program.Instructions -}}
	{{- with add $i.Op }}
	{{ $i.Output }}.Mul({{ .X }}, {{ .Y }})
	{{- end -}}
	{{- with double $i.Op }}
	{{ $i.Output }}.Square({{ .X }})
	{{- end -}}
	{{- with shift $i.Op -}}
	{{- $first := 0 -}}
	{{- if ne $i.Output.Identifier .X.Identifier }}
	{{ $i.Output }}.Square({{ .X }})
	{{- $first = 1 -}}
	{{- end }}
	for s := {{ $first }}; s < {{ .S }}; s++ {
		{{ $i.Output }}.Square({{ $i.Output }})
	}
	{{- end -}}
	{{- end }}
	return e.Set(z)
}
`

const tmplAddchainExp2 = `
func (e *Element) expP(x *Element) *Element {
	// The sequence of {{ .Ops.Adds }} multiplications and {{ .Ops.Doubles }} squarings is derived from the
	// following addition chain generated with {{ .Meta.Module }} {{ .Meta.ReleaseTag }}.
	//
	{{- range lines (format .Script) }}
	//	{{ . }}
	{{- end }}
	//
	var z = new(Element).Set(e)
	{{- range .Program.Temporaries }}
	var {{ . }} = new(Element)
	{{- end }}
	{{ range $i := .Program.Instructions -}}
	{{- with add $i.Op }}
	{{ $i.Output }}.Mul({{ .X }}, {{ .Y }})
	{{- end -}}
	{{- with double $i.Op }}
	{{ $i.Output }}.Square({{ .X }})
	{{- end -}}
	{{- with shift $i.Op -}}
	{{- $first := 0 -}}
	{{- if ne $i.Output.Identifier .X.Identifier }}
	{{ $i.Output }}.Square({{ .X }})
	{{- $first = 1 -}}
	{{- end }}
	for s := {{ $first }}; s < {{ .S }}; s++ {
		{{ $i.Output }}.Square({{ $i.Output }})
	}
	{{- end -}}
	{{- end }}
	return e.Set(z)
}
`

const tmplAddchainExp3 = `
func (e *Element) expPMinus1Over2(x *Element) *Element {
	// The sequence of {{ .Ops.Adds }} multiplications and {{ .Ops.Doubles }} squarings is derived from the
	// following addition chain generated with {{ .Meta.Module }} {{ .Meta.ReleaseTag }}.
	//
	{{- range lines (format .Script) }}
	//	{{ . }}
	{{- end }}
	//
	var z = new(Element).Set(e)
	{{- range .Program.Temporaries }}
	var {{ . }} = new(Element)
	{{- end }}
	{{ range $i := .Program.Instructions -}}
	{{- with add $i.Op }}
	{{ $i.Output }}.Mul({{ .X }}, {{ .Y }})
	{{- end -}}
	{{- with double $i.Op }}
	{{ $i.Output }}.Square({{ .X }})
	{{- end -}}
	{{- with shift $i.Op -}}
	{{- $first := 0 -}}
	{{- if ne $i.Output.Identifier .X.Identifier }}
	{{ $i.Output }}.Square({{ .X }})
	{{- $first = 1 -}}
	{{- end }}
	for s := {{ $first }}; s < {{ .S }}; s++ {
		{{ $i.Output }}.Square({{ $i.Output }})
	}
	{{- end -}}
	{{- end }}
	return e.Set(z)
}
`

const tmplAddchainInvert = `// Code generated by {{ .Meta.Name }}. DO NOT EDIT.
package bn256
// Invert sets e = 1/x, and returns e.
//
// If x == 0, Invert returns e = 0.
func (e *Element) Invert(x *Element) *Element {
	// Inversion is implemented as exponentiation with exponent p − 2.
	// The sequence of {{ .Ops.Adds }} multiplications and {{ .Ops.Doubles }} squarings is derived from the
	// following addition chain generated with {{ .Meta.Module }} {{ .Meta.ReleaseTag }}.
	//
	{{- range lines (format .Script) }}
	//	{{ . }}
	{{- end }}
	//
	var z = new(Element).Set(e)
	{{- range .Program.Temporaries }}
	var {{ . }} = new(Element)
	{{- end }}
	{{ range $i := .Program.Instructions -}}
	{{- with add $i.Op }}
	{{ $i.Output }}.Mul({{ .X }}, {{ .Y }})
	{{- end -}}
	{{- with double $i.Op }}
	{{ $i.Output }}.Square({{ .X }})
	{{- end -}}
	{{- with shift $i.Op -}}
	{{- $first := 0 -}}
	{{- if ne $i.Output.Identifier .X.Identifier }}
	{{ $i.Output }}.Square({{ .X }})
	{{- $first = 1 -}}
	{{- end }}
	for s := {{ $first }}; s < {{ .S }}; s++ {
		{{ $i.Output }}.Square({{ $i.Output }})
	}
	{{- end -}}
	{{- end }}
	return e.Set(z)
}
`

const tmplAddchainSqrt = `
// Sqrt sets e to a square root of x. If x is not a square, Sqrt returns
// false and e is unchanged. e and x can overlap.
func Sqrt(e, x *Element) (isSquare bool) {
	candidate, b, i := &gfP{}, &gfP{}, &gfP{}
	sqrtCandidate(candidate, x)
	gfpMul(b, twoExpPMinus5Over8, candidate) // b=ta1
	gfpMul(candidate, x, b)                  // a1=fb
	gfpMul(i, two, candidate)                // i=2(fb)
	gfpMul(i, i, b)                   // i=2(fb)b
	gfpSub(i, i, one)                 // i=2(fb)b-1
	gfpMul(i, candidate, i)                  // i=(fb)(2(fb)b-1)
	square := new(Element).Square(i)
	if square.Equal(x) != 1 {
		return false
	}
	e.Set(i)
	return true
}

// sqrtCandidate sets z to a square root candidate for x. z and x must not overlap.
func sqrtCandidate(z, x *Element) {
	// Since p = 8k+5, exponentiation by (p - 5) / 8 yields a square root candidate.
	//
	// The sequence of {{ .Ops.Adds }} multiplications and {{ .Ops.Doubles }} squarings is derived from the
	// following addition chain generated with {{ .Meta.Module }} {{ .Meta.ReleaseTag }}.
	//
	{{- range lines (format .Script) }}
	//	{{ . }}
	{{- end }}
	//
	{{- range .Program.Temporaries }}
	var {{ . }} = new(Element)
	{{- end }}
	{{ range $i := .Program.Instructions -}}
	{{- with add $i.Op }}
	{{ $i.Output }}.Mul({{ .X }}, {{ .Y }})
	{{- end -}}
	{{- with double $i.Op }}
	{{ $i.Output }}.Square({{ .X }})
	{{- end -}}
	{{- with shift $i.Op -}}
	{{- $first := 0 -}}
	{{- if ne $i.Output.Identifier .X.Identifier }}
	{{ $i.Output }}.Square({{ .X }})
	{{- $first = 1 -}}
	{{- end }}
	for s := {{ $first }}; s < {{ .S }}; s++ {
		{{ $i.Output }}.Square({{ $i.Output }})
	}
	{{- end -}}
	{{- end }}
}
`
