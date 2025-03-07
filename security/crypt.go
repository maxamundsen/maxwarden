package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"maxwarden/config"
	"os"

	"github.com/btcsuite/btcutil/base58"

	"github.com/minio/highwayhash"
	"golang.org/x/crypto/bcrypt"
)

////////////////////////////////
// Encoding Wrappers
////////////////////////////////

func EncodeBase64(in []byte) string {
	return base64.StdEncoding.EncodeToString(in)
}

func DecodeBase64(in string) []byte {
	out, _ := base64.RawStdEncoding.DecodeString(in)
	return out
}

func EncodeBase58(in []byte) string {
	return base58.Encode(in)
}

func DecodeBase58(in string) []byte {
	return base58.Decode(in)
}

////////////////////////////////
// HASH FUNCTIONS
////////////////////////////////

// Hash with SHA512 and output a Base58 string
func SHA512_58(in string) string {
	hasher := sha512.New()

	hasher.Write([]byte(in))

	hashBytes := hasher.Sum(nil)

	hashString := base58.Encode(hashBytes)

	return hashString
}

func HighwayHash58(in string) (string, error) {
	key := []byte(config.DATA_HASH_KEY)

	hasher, err := highwayhash.New(key)
	if err != nil {
		log.Println("Error generating hasher.")
		return "", err
	}

	hasher.Write([]byte(in))

	hash := hasher.Sum(nil)

	encodedData := base58.Encode(hash)

	return encodedData, nil
}

func HighwayHash(in string) (string, error) {
	key := []byte(config.DATA_HASH_KEY)

	hasher, err := highwayhash.New(key)
	if err != nil {
		log.Println("Error generating hasher.")
		return "", err
	}

	hasher.Write([]byte(in))

	hash := hasher.Sum(nil)

	return base64.StdEncoding.EncodeToString(hash), nil
}

func QuickFileHash(filepath string) (string, error) {
	key := []byte(config.DATA_HASH_KEY)

	file, err := os.Open(filepath)
	if err != nil {
		log.Println("Error opening file", err)
		return "", err
	}
	defer file.Close()

	hasher, err := highwayhash.New(key)
	if err != nil {
		log.Println("Error generating hasher.")
		return "", err
	}

	_, err = io.Copy(hasher, file)
	if err != nil {
		log.Println("Error hashing file:", err)
		return "", err
	}

	hash := hasher.Sum(nil)

	return base64.StdEncoding.EncodeToString(hash), nil
}

// Hash password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compare password with hash using bcrypt
func ComparePasswords(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RandBase58String(entropyBytes int) string {
	b := make([]byte, entropyBytes)
	rand.Read(b)
	return base58.Encode(b)
}

////////////////////////////////
// Encryption FUNCTIONS
////////////////////////////////


// AES Encrypt
func EncryptSecret(data []byte, passKey string) ([]byte, error) {
	key := make([]byte, 32)
	copy(key, passKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encryptedData := gcm.Seal(nonce, nonce, data, nil)

	return encryptedData, nil
}

// AES Decrypt
func DecryptSecret(encryptedData []byte, passKey string) ([]byte, error) {
	key := make([]byte, 32)
	copy(key, passKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, encryptedData := encryptedData[:nonceSize], encryptedData[nonceSize:]

	decryptedData, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

// Encrypt data using default private key
func EncryptData[T any](data *T) ([]byte, error) {
	return EncryptDataWithKey(data, config.GetConfig().IdentityPrivateKey)
}

// Decrypt data using default private key
func DecryptData[T any](data []byte) (*T, error) {
	return DecryptDataWithKey[T](data, config.GetConfig().IdentityPrivateKey)
}

func EncryptDataWithKey[T any](data *T, key string) ([]byte, error) {
	// serialize
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(data)
	if err != nil {
		return nil, err
	}

	// encrypt
	out, err := EncryptSecret(b.Bytes(), key)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func DecryptDataWithKey[T any](data []byte, key string) (*T, error) {
	dest := new(T)

	secret, err := DecryptSecret(data, key)
	if err != nil {
		return nil, err
	}

	// de-serialized
	b := bytes.Buffer{}
	b.Write(secret)

	d := gob.NewDecoder(&b)
	gobErr := d.Decode(dest)
	if gobErr != nil {
		return nil, gobErr
	}

	return dest, nil
}
