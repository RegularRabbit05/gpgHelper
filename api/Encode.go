package api

import (
	"io"
	"net/http"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

var pgpDecodeInstance = crypto.PGP()

func EncodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	toEncode := r.Header.Get("Key-Payload")
	if toEncode == "" {
		http.Error(w, "Missing Key-Payload header", http.StatusBadRequest)
		return
	}

	pubKey, err := io.ReadAll(r.Body)
	if err != nil || len(pubKey) == 0 {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	publicKey, err := crypto.NewKeyFromArmored(string(pubKey))
	if err != nil {
		http.Error(w, "Invalid public key", http.StatusBadRequest)
		return
	}

	encHandle, err := pgpDecodeInstance.Encryption().Recipient(publicKey).New()
	if err != nil {
		http.Error(w, "Failed to create encryption handle", http.StatusInternalServerError)
		return
	}

	pgpMessage, err := encHandle.Encrypt([]byte(toEncode))
	if err != nil {
		http.Error(w, "Failed to encrypt message", http.StatusInternalServerError)
		return
	}

	armoredMessage, err := pgpMessage.ArmorBytes()
	if err != nil {
		http.Error(w, "Failed to armor message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(armoredMessage)
}
