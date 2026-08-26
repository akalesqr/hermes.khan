// gen-license generates Ed25519 keys + signed license files for Khan.
// Usage:
//
//	gen-license genkeys           → writes private.key + public.key
//	gen-license sign <priv.key> <company> <max_users> <expiry YYYY-MM-DD> [-out license.key]
//
// The private key NEVER ships with the product. Only public.key is embedded.
package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"khan/internal/service"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		usage()
		return
	}

	switch args[0] {
	case "genkeys":
		genKeys()
	case "sign":
		sign(args[1:])
	case "pubkey-go":
		pubkeyGo(args[1:])
	default:
		usage()
	}
}

func usage() {
	fmt.Print(`Khan license tool

Commands:
  genkeys                        generate private.key + public.key (Ed25519)
  sign <priv.key> <company> <max_users> <expiry> [-out license.key]
                                 create a signed license (expiry: YYYY-MM-DD)
  pubkey-go <public.key>         print Go source for embedding the public key
`)
}

func genKeys() {
	priv, pub, err := service.GenerateKeyPair()
	if err != nil {
		fatal(err)
	}
	privB64 := base64.StdEncoding.EncodeToString(priv.Seed())
	pubB64 := base64.StdEncoding.EncodeToString(pub)
	if err := os.WriteFile("private.key", []byte(privB64+"\n"), 0o600); err != nil {
		fatal(err)
	}
	if err := os.WriteFile("public.key", []byte(pubB64+"\n"), 0o644); err != nil {
		fatal(err)
	}
	fmt.Println("✅ private.key written (KEEP SECRET — never distribute)")
	fmt.Println("✅ public.key written (embed in binary)")
	fmt.Println("Public key (base64):", pubB64)
}

func sign(args []string) {
	if len(args) < 4 {
		usage()
		os.Exit(1)
	}
	privPath, company := args[0], args[1]
	var maxUsers int
	var expiryStr string
	out := "license.key"

	// parse maxUsers and expiry from remaining args (support -out flag)
	rest := args[2:]
	for i := 0; i < len(rest); i++ {
		if rest[i] == "-out" && i+1 < len(rest) {
			out = rest[i+1]
			i++
			continue
		}
		if maxUsers == 0 {
			fmt.Sscanf(rest[i], "%d", &maxUsers)
			continue
		}
		expiryStr = rest[i]
	}
	if maxUsers == 0 {
		fmt.Sscanf(args[2], "%d", &maxUsers)
		expiryStr = args[3]
	}

	seedB64, err := os.ReadFile(privPath)
	if err != nil {
		fatal(err)
	}
	seed, err := base64.StdEncoding.DecodeString(string(seedB64))
	if err != nil {
		fatal(err)
	}
	priv := ed25519.NewKeyFromSeed(seed)

	exp, err := time.Parse("2006-01-02", expiryStr)
	if err != nil {
		fatal("bad expiry date, use YYYY-MM-DD")
	}

	lic, err := service.SignLicense(priv, company, maxUsers, exp)
	if err != nil {
		fatal(err)
	}
	if err := os.WriteFile(out, lic, 0o644); err != nil {
		fatal(err)
	}
	fmt.Printf("✅ license.key written → %s | %d users | expires %s\n", company, maxUsers, exp.Format("2006-01-02"))
}

func pubkeyGo(args []string) {
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		fatal(err)
	}
	pubB64 := string(data)
	_, err = base64.StdEncoding.DecodeString(pubB64)
	if err != nil {
		fatal(err)
	}
	fmt.Println("// Public key for Khan license verification (auto-generated)")
	fmt.Println("// Place in internal/service/pubkey.go")
	fmt.Printf("var khanPubKeyB64 = %q\n\n", pubB64)
	fmt.Printf("func init() { pk, err := base64.StdEncoding.DecodeString(khanPubKeyB64); if err == nil { SetPublicKey(ed25519.PublicKey(pk)) } }\n")
	_ = hex.EncodeToString
	_ = json.Marshal
}

func fatal(v any) {
	fmt.Fprintln(os.Stderr, "❌", v)
	os.Exit(1)
}
