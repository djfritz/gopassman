package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
)

const RANDOM_KEY_LENGTH = 16
const SALT_LENGTH = 128
const BUFFER_SIZE = 1024

func encrypt(password string, keys keymap) {
	log_info("encrypting %v keys", len(keys))

	// add to the password a 128 byte salt
	salt := generate_salt()
	log_info("generated salt:", salt)

	// take the sha256 of the password for the aes256 key
	h := sha256.New()
	key := make([]byte, sha256.Size)
	h.Write([]byte(password + salt))
	h.Sum(key[:0])
	log_info("sha256:", key)

	// generate the cipher for aes256
	block, err := aes.NewCipher(key)
	if err != nil {
		log_dieln(err)
	}

	// generate a checksum of the plaintext, and append
	plaintext := string(keys.Encode())
	h.Reset()
	checksum := make([]byte, sha256.Size)
	h.Write([]byte(plaintext))
	h.Sum(checksum[:0])
	log_info("plaintext checksum:", checksum)
	plaintext += string(checksum)

	// encrypt the plaintext keys + checksum
	ciphertext := crypto_encrypt(block, []byte(plaintext))
	log_info("ciphertext is %v bytes\n", len(ciphertext))

	// generate another checksum for integrity
	h.Reset()
	integrity := make([]byte, sha256.Size)
	h.Write(ciphertext)
	h.Sum(integrity[:0])
	log_info("integrity:", integrity)

	// save to file, along with salt, and integrity
	file := file_create()
	file.Write(integrity)
	file.Write([]byte(salt))
	file.Write(ciphertext)
	file.Close()
}

func decrypt(file *os.File, password string) keymap {
	keys := make(keymap)

	// read everything in
	integrity := make([]byte, sha256.Size)
	_, err := file.Read(integrity[:])
	if err != nil {
		log_dieln(err)
	}
	log_info("read integrity:", integrity)
	salt := make([]byte, SALT_LENGTH)
	_, err = file.Read(salt[:])
	if err != nil {
		log_dieln(err)
	}
	log_info("read salt:", string(salt))
	var ciphertext string
	buffer := make([]byte, BUFFER_SIZE)
	for {
		n, err := file.Read(buffer[:])
		ciphertext += string(buffer[:n])
		if err == io.EOF {
			break
		} else if err != nil {
			log_dieln(err)
		}
	}
	log_info("read %v bytes of ciphertext\n", len(ciphertext))

	// check integrity
	h := sha256.New()
	integrity_check := make([]byte, sha256.Size)
	h.Write([]byte(ciphertext))
	h.Sum(integrity_check[:0])
	log_info("integrity check:", integrity_check)
	if string(integrity) != string(integrity_check) {
		log_info("original integrity:", integrity)
		log_dieln("file integrity error!")
	}

	// take the sha256 of the password for the aes256 key
	h.Reset()
	key := make([]byte, sha256.Size)
	h.Write([]byte(password + string(salt)))
	h.Sum(key[:0])
	log_info("sha256:", key)

	// generate the cipher for aes256
	block, err := aes.NewCipher(key)
	if err != nil {
		log_dieln(err)
	}

	// decrypt using ciphertext
	ptcs := crypto_decrypt(block, []byte(ciphertext))

	// read checksum and plaintext from ptcs
	plaintext := ptcs[:len(ptcs)-sha256.Size]
	checksum := ptcs[len(ptcs)-sha256.Size:]

	log_info("got %v bytes of plaintext\n", len(plaintext))
	log_infoln("checksum:", checksum)

	// take checksum of plaintext and compare with checksum
	h.Reset()
	checksum_check := make([]byte, sha256.Size)
	h.Write(plaintext)
	h.Sum(checksum_check[:0])
	log_info("checksum check:", checksum_check)
	if string(checksum) != string(checksum_check) {
		log_dieln("invalid password")
	}

	// build keys from plaintext, return
	keys.Decode(plaintext)

	return keys
}

func generate_salt() string {
	return string(generate_random_bytes(SALT_LENGTH))
}

func crypto_encrypt(block cipher.Block, value []byte) []byte {
	// generate an IV
	// http://en.wikipedia.org/wiki/Block_cipher_modes_of_operation
	iv := make([]byte, block.BlockSize())
	rand.Read(iv)
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(value, value)
	return append(iv, value...)
}

func crypto_decrypt(block cipher.Block, value []byte) []byte {
	if len(value) > block.BlockSize() {
		// Extract iv.
		iv := value[:block.BlockSize()]
		// Extract ciphertext.
		value = value[block.BlockSize():]
		// Decrypt it.
		stream := cipher.NewCTR(block, iv)
		stream.XORKeyStream(value, value)
		return value
	}
	return nil
}

// generate a url encodable base64 string as a random password
func generate_key() []byte {
	return generate_random_bytes(RANDOM_KEY_LENGTH)
}

func generate_random_bytes(l int) []byte {
	b := make([]byte, l)
	rand.Read(b)
	e := base64.URLEncoding
	d := make([]byte, e.EncodedLen(len(b)))
	e.Encode(d, b)
	return d[:l]
}
