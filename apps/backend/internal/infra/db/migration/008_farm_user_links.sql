-- CafeOS - Migration 008: Producer vira vínculo N:N Farm↔User com papel
-- NOTE: schema is authoritatively managed by GORM AutoMigrate (see
-- internal/infra/db/postgres/connection.go). This file exists for parity
-- with the docker-compose bootstrap and as documentation of the change;
-- it is not run by a migration tool.
--
-- producers deixa de ser 1:1 com farms (era reforçado por
-- UNIQUE(farm_id) em 002_farm_producer_fields.sql) e vira o vínculo N:N
-- entre uma fazenda e os usuários que atuam nela — cada linha é "este User
-- tem este papel (role_id) nesta fazenda" (ex: proprietário, operador de
-- campo, financeiro). Uma fazenda pode ter vários vínculos; um usuário
-- pode estar vinculado a várias fazendas. RoleID reaproveita o catálogo
-- global de roles (ver 007_roles_modules.sql) — é só descritivo aqui, não
-- altera o papel/permissão global do usuário no sistema.
--
-- Os campos cadastrais (cpf, rg, etc.) continuam existindo por
-- compatibilidade com o registro legal do produtor (CAR/INCRA), relevantes
-- sobretudo quando o vínculo é do papel "proprietario" — nada no banco
-- força isso, é convenção de uso na tela.

ALTER TABLE producers DROP CONSTRAINT IF EXISTS producers_farm_id_key;
DROP INDEX IF EXISTS idx_producers_farm_id_key;

ALTER TABLE producers ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE producers ADD COLUMN IF NOT EXISTS role_id UUID NOT NULL;

ALTER TABLE producers
    ADD CONSTRAINT fk_producers_role FOREIGN KEY (role_id) REFERENCES roles(id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_producers_farm_user
    ON producers (farm_id, user_id);
