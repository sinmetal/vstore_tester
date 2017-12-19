package item

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"golang.org/x/oauth2/google"

	"google.golang.org/api/cloudkms/v1"

	"github.com/pkg/errors"

	"github.com/sinmetal/vstore_tester/config"
)

var kmsClient *http.Client
var kmsService *cloudkms.Service

func SetUpKMS(ctx context.Context) error {
	var err error
	kmsClient, err = google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		return errors.Wrap(err, "encrypt: failed google.DefaultClient: ")
	}

	// Create the KMS client.
	kmsService, err = cloudkms.New(kmsClient)
	if err != nil {
		return errors.Wrap(err, "encrypt: failed cloudkms.New: ")
	}

	return nil
}

func encrypt(ctx context.Context, plaintext string) (ciphertext string, cryptoKey string, err error) {
	projectID, err := config.GetProjectID(ctx)
	if err != nil {
		return "", "", err
	}

	cryptoKeyName := fmt.Sprintf("projects/%s/locations/global/keyRings/sample-ring/cryptoKeys/sample-key", projectID)
	response, err := kmsService.Projects.Locations.KeyRings.CryptoKeys.Encrypt(cryptoKeyName, &cloudkms.EncryptRequest{
		Plaintext: base64.StdEncoding.EncodeToString([]byte(plaintext)),
	}).Do()
	if err != nil {
		return "", "", errors.Wrapf(err, "encrypt: failed to encrypt. CryptoKey=%s", cryptoKeyName)
	}

	return response.Ciphertext, response.Name, nil
}

func decrypt(ctx context.Context, ciphertext string) (plaintext string, err error) {
	projectID, err := config.GetProjectID(ctx)
	if err != nil {
		return "", err
	}

	cryptoKeyName := fmt.Sprintf("projects/%s/locations/global/keyRings/sample-ring/cryptoKeys/sample-key", projectID)

	response, err := kmsService.Projects.Locations.KeyRings.CryptoKeys.Decrypt(cryptoKeyName, &cloudkms.DecryptRequest{
		Ciphertext: ciphertext,
	}).Do()
	if err != nil {
		return "", errors.Wrapf(err, "decrypt: failed to decrypt. CryptoKey=%s", cryptoKeyName)
	}

	return response.Plaintext, nil
}
