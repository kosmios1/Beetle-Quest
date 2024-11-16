package service

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/repositories"
)

type AdminService struct {
	arepo repositories.AdminRepo
}

func NewAdminService(arepo repositories.AdminRepo) *AdminService {
	return &AdminService{arepo: arepo}
}

func (s *AdminService) FindByID(aid models.UUID) (*models.Admin, bool) {
	return s.arepo.FindByID(aid)
}
