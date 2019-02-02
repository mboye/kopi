package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"hash"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/pbkdf2"
)

const (
	keyIterations  = 100000
	passwordEnvVar = "KOPI_PASSWORD"
	keyLength      = 32  // Bytes for AES-256,
	SaltLength     = 128 // Bytes
	SaltFileName   = "salt"
)

type Context struct {
	Salt        []byte
	newHash     func() hash.Hash
	CipherBlock cipher.Block
}

func NewContext(outputDir string, withCrypto bool) (*Context, error) {
	ctx := &Context{
		newHash: sha1.New}

	if err := ctx.loadSalt(outputDir); err != nil {
		return nil, err
	}

	if withCrypto {
		if err := ctx.initCrypto(); err != nil {
			return nil, err
		}
	}

	return ctx, nil
}

func (ctx *Context) loadSalt(outputDir string) error {
	saltPath := fmt.Sprintf("%s/%s", outputDir, SaltFileName)
	ctx.Salt = make([]byte, SaltLength)
	if saltFile, err := os.Open(saltPath); err == nil {
		// Read existing salt
		bytesRead, err := saltFile.Read(ctx.Salt)
		if err != nil {
			return fmt.Errorf("failed to read salt: %s", err.Error())
		} else if bytesRead != SaltLength {
			return fmt.Errorf("incomplete salt read: %d of %d bytes read", bytesRead, SaltLength)
		}
	} else if os.IsNotExist(err) {
		// Create salt
		if bytesRead, err := rand.Read(ctx.Salt); err != nil || bytesRead != SaltLength {
			return fmt.Errorf("failed to generate salt: %s", err.Error())
		} else {
			if err := ioutil.WriteFile(saltPath, ctx.Salt, 0655); err != nil {
				return fmt.Errorf("failed to save salt: %s", err.Error())
			}
			log.Info("salt created")
		}
	} else {
		return fmt.Errorf("failed to open salt file: %s", err.Error())
	}
	return nil
}

func (ctx *Context) initCrypto() error {
	if passwordString, found := os.LookupEnv(passwordEnvVar); !found {
		return fmt.Errorf("environment variable undefined: %s", passwordEnvVar)
	} else {
		password := []byte(passwordString)
		key := pbkdf2.Key(password, ctx.Salt, keyIterations, keyLength, sha1.New)
		block, err := aes.NewCipher(key)
		if err != nil {
			return fmt.Errorf("failed to init cipher: %s", err.Error())
		}
		ctx.CipherBlock = block
		return nil
	}
}

func (ctx *Context) NewHasher() (hash.Hash, error) {
	hasher := ctx.newHash()
	if bytesWritten, err := hasher.Write(ctx.Salt); err != nil {
		return nil, fmt.Errorf("failed to write salt: %s", err.Error())
	} else if bytesWritten != SaltLength {
		return nil, fmt.Errorf("incomplete salt write: %d of %d bytes written", bytesWritten, SaltLength)
	}

	return hasher, nil
}
