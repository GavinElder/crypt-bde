package bde

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

// GetActiveKeyProtector returns the primary activation key (not the TPM)
func GetActiveKeyProtector() (string, error) {
	keys, err := GetKeyProtectors()
	if err != nil {
		return "", errors.Wrap(err, "getting active key protector")
	}

	var kp []string

	for _, key := range keys {
		if key.KeyProtectorType == 3 {
			kp = append(kp, key.KeyProtectorID)
		}
	}

	ak := kp[0]

	return ak, nil
}

// GetKeyProtectors Lists All Active Key Protectors on the System Drive
func GetKeyProtectors() ([]KeyProtectors, error) {
	cmd := exec.Command("powershell", "(Get-BitlockerVolume -MountPoint $env:SystemDrive).KeyProtector", "|", "ConvertTo-Json")

	// cmd.Stderr = os.Stderr
	o, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "exec Get-BitlockerVolume")
	}

	var kp []KeyProtectors

	if err := json.Unmarshal(o, &kp); err != nil {
		return nil, errors.Wrap(err, "failed unmarshalling Key Protectors")
	}

	fmt.Println(kp)

	return kp, nil
}

// KeyProtectors represent each item that can unlock the disk
type KeyProtectors struct {
	KeyProtectorID   string `json:"KeyProtectorId"`
	KeyProtectorType int    `json:"KeyProtectorType"`
	RecoveryPassword string `json:"RecoveryPassword"`
}
