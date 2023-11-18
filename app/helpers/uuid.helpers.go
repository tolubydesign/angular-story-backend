package helpers

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/google/uuid"
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

func GenerateStringUUID() string {
	str := uuid.New().String()
	fmt.Println("id:", str)
	return str
}
