package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"encoding/hex"
	"log"
)

type AdminRepo struct {
	admins map[models.UUID]models.Admin
}

func NewAdminRepo() *AdminRepo {
	repo := &AdminRepo{
		admins: make(map[models.UUID]models.Admin),
	}

	uuid, err := utils.ParseUUID("09087f45-5209-4efa-85bd-761562a6df53")
	if err != nil {
		log.Fatalf("Invalid admin's UUID in mock admin repo\n")
	}

	passHash, err := hex.DecodeString("243261243130247370373732344b6d544a302e4f347862557176514d754c5330464a79684e4355736c6e59787757685a6668386a7739704430644457")
	if err != nil {
		log.Fatalf("Invalid admin's password hash \n")
	}

	repo.admins[uuid] = models.Admin{
		AdminId:      uuid,
		PasswordHash: passHash,
		OtpSecret:    "g2ytwh764px5wzorxcbk2c2f2jhv74kd",
		Email:        "admin@admin.com",
	}

	return repo
}

func (r *AdminRepo) FindByID(id models.UUID) (*models.Admin, error) {
	if admin, ok := r.admins[id]; ok {
		return &admin, nil
	}
	return nil, models.ErrAdminNotFound
}
