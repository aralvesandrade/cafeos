-- CafeOS - Migration 009: hierarquia de usuários (principal/sub-usuário) e
-- plano movido de organização para usuário.
-- NOTE: schema is authoritatively managed by GORM AutoMigrate (see
-- internal/infra/db/postgres/connection.go). This file exists for parity
-- with the docker-compose bootstrap and as documentation of the change;
-- it is not run by a migration tool.
--
-- Organization já funciona como tenant/whitelabel (CafeOS é o tenant padrão
-- da plataforma; um whitelabel tipo "Cooperativa XPTO" é outro tenant).
-- Dentro de um mesmo tenant pode haver vários usuários principais
-- independentes, cada um dono de um grupo de fazendas e responsável pelo
-- próprio Plan. users.managed_by_user_id (auto-FK, nulo = principal) marca
-- essa hierarquia de 2 níveis: um principal cria sub-usuários (que recebem
-- managed_by_user_id apontando para ele); sub-usuário nunca cria outro
-- usuário nem tem plan_id próprio.

ALTER TABLE users ADD COLUMN IF NOT EXISTS managed_by_user_id UUID REFERENCES users(id);
ALTER TABLE users ADD COLUMN IF NOT EXISTS plan_id UUID REFERENCES plans(id);

-- Backfill (idempotente): managed_by_user_id fica NULL para todos os
-- usuários existentes (todo mundo vira principal por padrão — ninguém
-- perde acesso). O plan_id do usuário mais antigo de cada organização
-- recebe o plan_id (antigo) da própria organização.
WITH oldest_user_per_org AS (
    SELECT DISTINCT ON (organization_id) id, organization_id
    FROM users
    ORDER BY organization_id, created_at ASC
)
UPDATE users u
SET plan_id = o.plan_id
FROM oldest_user_per_org ou
JOIN organizations o ON o.id = ou.organization_id
WHERE u.id = ou.id
  AND o.plan_id IS NOT NULL
  AND u.plan_id IS NULL;

-- Passo final opcional, DESTRUTIVO — só rodar após confirmar que nenhum
-- código ainda lê organizations.plan/organizations.plan_id:
-- ALTER TABLE organizations DROP COLUMN IF EXISTS plan;
-- ALTER TABLE organizations DROP COLUMN IF EXISTS plan_id;
