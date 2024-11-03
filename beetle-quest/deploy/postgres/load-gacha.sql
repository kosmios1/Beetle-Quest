CREATE TYPE rarity AS ENUM (
    'common', 'uncommon', 'rare', 'epic', 'legendary'
);

CREATE TABLE gachas (
    gacha_id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    rarity rarity NOT NULL,
    price BIGINT NOT NULL,
    image_path TEXT
);

INSERT INTO gachas (gacha_id, name, rarity, price, image_path) VALUES
(uuid_generate_v4(), 'Mage Butterfly', 'epic', 10000, '/images/mage_butterfly_epic.webp'),
(uuid_generate_v4(), 'Warrior Dragonfly', 'epic', 10000, '/images/warrior_dragonfly_epic.webp'),
(uuid_generate_v4(), 'WarriorHercule Beetle', 'legendary', 30000, '/images/warrior_hercule_beetle_legendary.webp');
-- ('uuid-generate-v4()', 'Example Gacha 2', 'rare', 1500, '/images/example2.png'),
-- ('uuid-generate-v4()', 'Example Gacha 2', 'rare', 1500, '/images/example2.png'),
-- ('uuid-generate-v4()', 'Example Gacha 2', 'rare', 1500, '/images/example2.png'),
-- ('uuid-generate-v4()', 'Example Gacha 2', 'rare', 1500, '/images/example2.png'),
