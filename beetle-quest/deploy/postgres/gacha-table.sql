CREATE TABLE gachas (
    gacha_id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    rarity INTEGER NOT NULL,
    price BIGINT NOT NULL,
    image_path TEXT
);

INSERT INTO gachas (gacha_id, name, rarity, price, image_path) VALUES
(uuid_generate_v4(), 'Warrior Locust',     0, 3000, '/static/images/warrior_locust_common.webp'),
(uuid_generate_v4(), 'Warrior Cricket',    0, 3000, '/static/images/warrior_cricket_common.webp'),
(uuid_generate_v4(), 'Tank Mole Cricket',  0, 3000, '/static/images/tank_mole-cricket_common.webp'),
(uuid_generate_v4(), 'Munich Grasshopper', 0, 3000, '/static/images/munich_grasshopper_common.webp'),

(uuid_generate_v4(), 'Warrior Centipede', 1, 5000, '/static/images/warrior_centipede_uncommon.webp'),
(uuid_generate_v4(), 'Priest Cicada',     1, 5000, '/static/images/priest_cicada_uncommon.webp'),
(uuid_generate_v4(), 'Mage Mosquito',     1, 5000, '/static/images/mage_mosquito_uncommon.webp'),
(uuid_generate_v4(), 'Druid Bee',         1, 5000, '/static/images/druid_bee_uncommon.webp'),

(uuid_generate_v4(), 'Warrior Beetle',    2, 7000, '/static/images/warrior_beetle_rare.webp'),
(uuid_generate_v4(), 'Tank Bee 1',        2, 7000, '/static/images/tank_bee_rare.webp'),
(uuid_generate_v4(), 'Priest Moth',       2, 7000, '/static/images/priest_moth_rare.webp'),
(uuid_generate_v4(), 'Druid Butterfly 1', 2, 7000, '/static/images/druid_butterfly_rare.webp'),
(uuid_generate_v4(), 'Assassin Mosquito', 2, 7000, '/static/images/assassin_mosquito_rare.webp'),

(uuid_generate_v4(), 'Mage Moth',         3, 10000, '/static/images/mage_moth_epic.webp'),
(uuid_generate_v4(), 'Mage Butterfly 1',  3, 10000, '/static/images/mage_butterfly_epic.webp'),
(uuid_generate_v4(), 'Assassin Peacock',  3, 10000, '/static/images/assassin_peacock_epic.webp'),
(uuid_generate_v4(), 'Mage Butterfly',    3, 10000, '/static/images/mage_butterfly_epic.webp'),
(uuid_generate_v4(), 'Warrior Dragonfly', 3, 10000, '/static/images/warrior_dragonfly_epic.webp'),

(uuid_generate_v4(), 'Tank Bee 2',             4, 30000, '/static/images/tank_bee_legendary.webp'),
(uuid_generate_v4(), 'Mage Butterfly 2',       4, 30000, '/static/images/mage_butterfly_legendary.webp'),
(uuid_generate_v4(), 'Druid Butterfly',        4, 30000, '/static/images/druid_butterfly_legendary.webp'),
(uuid_generate_v4(), 'Demoniac Mosquito',      4, 30000, '/static/images/demoniac_mosquito_legendary.webp'),
(uuid_generate_v4(), 'Warrior Hercule Beetle', 4, 30000, '/static/images/warrior_hercule_beetle_legendary.webp');
