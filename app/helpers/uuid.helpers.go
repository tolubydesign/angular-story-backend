package helpers

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os/exec"
)

func GoogleGenerateUUID() string {
	uuid, _ := uuid.NewUUID()
	return fmt.Sprintf("%s", uuid)
}

func OSGenerateUUID() string {
	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", uuid)
	return fmt.Sprintf("%s", uuid)
}
