package repository

import (
	"beetle-quest/pkg/models"
	"beetle-quest/pkg/utils"
	"sync"
)

type GachaRepo struct {
	mux        sync.RWMutex
	gachas     map[models.UUID]models.Gacha
	userGachas map[models.UUID][]models.UUID
}

func NewGachaRepo() *GachaRepo {
	repo := &GachaRepo{
		mux:        sync.RWMutex{},
		gachas:     make(map[models.UUID]models.Gacha),
		userGachas: make(map[models.UUID][]models.UUID),
	}

	populateMockRepo(repo)

	return repo
}

func (r *GachaRepo) Create(gacha *models.Gacha) error {
	r.mux.RLock()

	for _, g := range r.gachas {
		if gacha.Name == g.Name {
			return models.ErrGachaAlreadyExists
		}
	}

	if _, ok := r.gachas[gacha.GachaID]; ok {
		return models.ErrInternalServerError
	}
	r.mux.RUnlock()

	r.mux.Lock()
	defer r.mux.Unlock()
	r.gachas[gacha.GachaID] = *gacha

	return nil
}

func (r *GachaRepo) Update(gacha *models.Gacha) error {
	r.mux.RLock()

	for _, g := range r.gachas {
		if gacha.Name == g.Name {
			return models.ErrGachaAlreadyExists
		}
	}
	r.mux.RUnlock()

	r.mux.Lock()
	defer r.mux.Unlock()
	if _, ok := r.gachas[gacha.GachaID]; !ok {
		return models.ErrGachaNotFound
	}
	r.gachas[gacha.GachaID] = *gacha
	return nil
}

func (r *GachaRepo) Delete(gacha *models.Gacha) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, ok := r.gachas[gacha.GachaID]; !ok {
		return models.ErrGachaNotFound
	}
	delete(r.gachas, gacha.GachaID)
	return nil
}

func (r *GachaRepo) GetAll() ([]models.Gacha, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	var values []models.Gacha
	for _, value := range r.gachas {
		values = append(values, value)
	}
	return values, nil
}

func (r *GachaRepo) FindByID(gid models.UUID) (*models.Gacha, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	if user, ok := r.gachas[gid]; ok {
		return &user, nil
	}
	return nil, models.ErrGachaNotFound
}

func (r *GachaRepo) AddGachaToUser(uid, gid models.UUID) error {
	r.mux.RLock()
	// if _, ok := r.userGachas[uid]; !ok {
	// 	return models.ErrInternalServerError
	// }

	gachas := r.userGachas[uid]
	for _, gachaID := range gachas {
		if gachaID == gid {
			return models.ErrUserAlreadyHasGacha
		}
	}
	r.mux.RUnlock()

	r.mux.Lock()
	defer r.mux.Unlock()

	if _, ok := r.userGachas[uid]; !ok {
		r.userGachas[uid] = make([]models.UUID, 0)
	}
	r.userGachas[uid] = append(r.userGachas[uid], r.gachas[gid].GachaID)
	return nil
}

func (r *GachaRepo) RemoveGachaFromUser(uid, gid models.UUID) error {
	r.mux.RLock()
	if _, ok := r.userGachas[uid]; !ok {
		return models.ErrRetalationGachaUserNotFound
	}

	own := false
	gachaPos := -1
	gachas := r.userGachas[uid]
	for j, gachaID := range gachas {
		if gachaID == gid {
			own = true
			gachaPos = j
			break
		}
	}
	r.mux.RUnlock()

	if !own {
		return models.ErrRetalationGachaUserNotFound
	}

	r.mux.Lock()
	defer r.mux.Unlock()
	r.userGachas[uid] = append(r.userGachas[uid][:gachaPos], r.userGachas[uid][gachaPos+1:]...)
	return nil
}

func (r *GachaRepo) RemoveUserGachas(uid models.UUID) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, ok := r.userGachas[uid]; !ok {
		return models.ErrRetalationGachaUserNotFound
	}
	delete(r.userGachas, uid)
	return nil
}

func (r *GachaRepo) GetUserGachas(uid models.UUID) ([]models.Gacha, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	if gachaIDs, ok := r.userGachas[uid]; ok {
		var values []models.Gacha
		for _, gid := range gachaIDs {
			if g, ok := r.gachas[gid]; ok {
				values = append(values, g)
			}
		}
		return values, nil
	}
	return nil, models.ErrUserNotFound
}

// Utils ===============================================================================================================

func populateMockRepo(repo *GachaRepo) {
	{ // Add Gachas
		gachas := []models.Gacha{
			{GachaID: utils.PanicIfError[models.UUID](utils.ParseUUID("f6975151-5554-4160-b3cd-2d5aab548ccf")), Name: "Tank Mole Cricket", Rarity: models.Rarity(0), Price: 3000, ImagePath: "/static/images/tank_mole-cricket_common.webp"},
			{GachaID: utils.PanicIfError[models.UUID](utils.ParseUUID("7219d1c4-d71b-4603-9bbc-ce902bb3d895")), Name: "Warrior Locust", Rarity: models.Rarity(0), Price: 3000, ImagePath: "/static/images/warrior_locust_common.webp"},
			{GachaID: utils.PanicIfError[models.UUID](utils.ParseUUID("2e80dca8-51f4-46c9-9546-81c7cf78a4f7")), Name: "Warrior Cricket", Rarity: models.Rarity(0), Price: 3000, ImagePath: "/static/images/warrior_cricket_common.webp"},
			{GachaID: utils.PanicIfError[models.UUID](utils.ParseUUID("5a518a08-9c1d-4a1b-a6bf-191b3eb83456")), Name: "Munich Grasshopper", Rarity: models.Rarity(0), Price: 3000, ImagePath: "/static/images/munich_grasshopper_common.webp"},
			{GachaID: utils.PanicIfError[models.UUID](utils.ParseUUID("4b00f833-ac14-4f54-9ad0-9527b96d42b5")), Name: "Warrior Centipede", Rarity: models.Rarity(1), Price: 5000, ImagePath: "/static/images/warrior_centipede_uncommon.webp"},
			{GachaID: utils.PanicIfError[models.UUID](utils.ParseUUID("ed46194c-08b3-45a7-b735-7168a49f43ac")), Name: "Priest Cicada", Rarity: models.Rarity(1), Price: 5000, ImagePath: "/static/images/priest_cicada_uncommon.webp"},
			{GachaID: utils.PanicIfError[models.UUID](utils.ParseUUID("96f8ced4-0305-43ad-9e52-779013fa8502")), Name: "Mage Mosquito;", Rarity: models.Rarity(1), Price: 5000, ImagePath: "/static/images/mage_mosquito_uncommon.webp"},
			{GachaID: utils.PanicIfError[models.UUID](utils.ParseUUID("6f5d6b25-4b87-44f8-8878-3f06e6c9239c")), Name: "Druid Bee", Rarity: models.Rarity(1), Price: 5000, ImagePath: "/static/images/druid_bee_uncommon.webp"},
			{GachaID: utils.PanicIfError[models.UUID](utils.ParseUUID("c6cae6d7-e6aa-49c2-a5a3-5b6b2f435e6f")), Name: "Warrior Beetle", Rarity: models.Rarity(2), Price: 7000, ImagePath: "/static/images/warrior_beetle_rare.webp"},
			{GachaID: utils.PanicIfError[models.UUID](utils.ParseUUID("eda3bd94-322a-44fb-8268-b596cacb77c3")), Name: "Tank Bee 1", Rarity: models.Rarity(2), Price: 7000, ImagePath: "/static/images/tank_bee_rare.webp"},
			{GachaID: utils.PanicIfError[models.UUID](utils.ParseUUID("05311293-f0bf-4bbd-a879-72183c4cdab8")), Name: "Priest Moth", Rarity: models.Rarity(2), Price: 7000, ImagePath: "/static/images/priest_moth_rare.webp"},
			{GachaID: utils.PanicIfError[models.UUID](utils.ParseUUID("e455113c-655c-478d-bd24-b2a59c11e1f3")), Name: "Druid Butterfly 1", Rarity: models.Rarity(2), Price: 7000, ImagePath: "/static/images/druid_butterfly_rare.webp"},
			{GachaID: utils.PanicIfError[models.UUID](utils.ParseUUID("deb88718-311a-4aa1-a9da-afe2a901d5f6")), Name: "Assassin Mosquito", Rarity: models.Rarity(2), Price: 7000, ImagePath: "/static/images/assassin_mosquito_rare.webp"},
		}

		for _, g := range gachas {
			repo.gachas[g.GachaID] = g
		}
	}

	bobUID := utils.PanicIfError[models.UUID](utils.ParseUUID("744a2f4d-a693-4352-916e-64f4ef94b709"))
	repo.userGachas[bobUID] = make([]models.UUID, 0)
	repo.userGachas[bobUID] = append(repo.userGachas[bobUID], utils.PanicIfError[models.UUID](utils.ParseUUID("96f8ced4-0305-43ad-9e52-779013fa8502")))

	// TODO: populate the userGachas map
}
