package main

import (
	"crypto/aes"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	InputDir  string
	OutputDir string
	AES_KEY   = "_netsyna_netmod_"
)

func init() {
	flag.StringVar(&InputDir, "input", "configs", "Input directory containing .nm files")
	flag.StringVar(&OutputDir, "output", "decrypt", "Directory to save decrypted files")
	flag.Parse()
	flag.Usage = func() {
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
	if err := os.MkdirAll("configs", 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll("outputs", 0755); err != nil {
		panic(err)
	}
	files, err := filepath.Glob(filepath.Join(InputDir, "*.nm"))
	if err != nil || len(files) == 0 {
		fmt.Printf("No .nm files found in %s folder\n", InputDir)
		return
	}

	if err := os.MkdirAll(OutputDir, os.ModePerm); err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	for _, file := range files {
		fmt.Println("Decrypting:", file)

		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", file, err)
			continue
		}

		content := strings.TrimSpace(string(data))

		if strings.HasPrefix(content, "nm-") {
			// find the protocol
			idx := strings.Index(content, "://")
			if idx == -1 {
				fmt.Printf("Invalid format in %s, missing ://\n", file)
				continue
			}

			protocol := content[3:idx] // skip "nm-"
			raw := content[idx+3:]     // everything after ://

			ciphertext, err := base64.StdEncoding.DecodeString(raw)
			if err != nil {
				fmt.Printf("Base64 decode error for %s: %v\n", file, err)
				continue
			}

			plaintext, err := decryptAESECB(ciphertext, []byte(AES_KEY))
			if err != nil {
				fmt.Printf("Decrypt error for %s: %v\n", file, err)
				continue
			}

			outputFile := filepath.Join(OutputDir, strings.TrimSuffix(filepath.Base(file), ".nm")+".txt")
			decryptedString := protocol + "://" + string(trimNullBytes(plaintext))
			if err := os.WriteFile(outputFile, []byte(decryptedString), 0644); err != nil {
				fmt.Printf("Error writing %s: %v\n", outputFile, err)
				continue
			}

			fmt.Println("Saved decrypted file:", outputFile)
		} else {
			ciphertext, err := base64.StdEncoding.DecodeString(content)
			if err != nil {
				fmt.Printf("Base64 decode error for %s: %v\n", file, err)
				continue
			}
			plaintext, err := decryptAESECB(ciphertext, []byte(AES_KEY))
			if err != nil {
				fmt.Printf("Decrypt error for %s: %v\n", file, err)
				continue
			}
			outputFile := filepath.Join(OutputDir, strings.TrimSuffix(filepath.Base(file), ".nm")+".txt")
			if err := os.WriteFile(outputFile, trimNullBytes(plaintext), 0644); err != nil {
				fmt.Printf("Error writing %s: %v\n", outputFile, err)
				continue
			}

			fmt.Println("Saved decrypted file:", outputFile)
		}
	}

	fmt.Println("All files processed.")
}

func decryptAESECB(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("ciphertext length not multiple of block size")
	}

	plaintext := make([]byte, len(ciphertext))
	bs := block.BlockSize()
	for start := 0; start < len(ciphertext); start += bs {
		block.Decrypt(plaintext[start:start+bs], ciphertext[start:start+bs])
	}
	return plaintext, nil
}
func trimNullBytes(data []byte) []byte {
	return []byte(strings.TrimRight(string(data), "\x00"))
}
