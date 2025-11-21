package main

import (
	"fmt"

	"internal/auth"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestFirst(t *testing.T) {

	ustring, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")
	tokenstring, err := auth.MakeJWT(ustring, "secretsalt", time.Hour)
	if err != nil {
		t.Error("Error creating JWT")
	}
	fmt.Printf("%v\n", tokenstring)

	uuid, err := auth.ValidateJWT(tokenstring, "secretsalt")
	if err != nil {
		t.Error("Error validating jwt")
	}
	fmt.Printf("%v\n", uuid.String())
}

func TestSecond(t *testing.T) {

	ustring, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")
	tokenstring, err := auth.MakeJWT(ustring, "secret", time.Hour)
	if err != nil {
		t.Error("Error creating JWT")
	}
	fmt.Printf("%v\n", tokenstring)

	uuid, err := auth.ValidateJWT(tokenstring, "secretsalt")
	if err == nil {
		t.Error("Error validating jwt")
	}
	fmt.Printf("%v\n", uuid.String())

}
