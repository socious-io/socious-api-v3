-- Add Offer and Mission to contract
ALTER TABLE contracts
    ADD COLUMN offer_id UUID,
    ADD COLUMN mission_id UUID,
    ADD CONSTRAINT fk_offer FOREIGN KEY (offer_id) REFERENCES offers(id) ON DELETE SET NULL,
    ADD CONSTRAINT fk_mission FOREIGN KEY (mission_id) REFERENCES missions(id) ON DELETE SET NULL;

ALTER TYPE contract_status ADD VALUE 'COMPLETED';

-- Add triggers
CREATE OR REPLACE FUNCTION upsert_contract_on_offer() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    contract RECORD;
    project RECORD;
    contract_status TEXT;
    v_commitment_period TEXT;
    v_commitment_period_count INTEGER;
BEGIN
    SELECT * INTO contract
    FROM contracts
    WHERE offer_id = NEW.id;

    SELECT * INTO project
    FROM projects
    WHERE id = NEW.project_id AND payment_type IS NOT NULL;

    contract_status := CASE
        WHEN NEW.status='PENDING' THEN 'CREATED'
        WHEN NEW.status='APPROVED' THEN 'CLIENT_APPROVED'
        WHEN NEW.status='HIRED' THEN 'SIGNED'
        WHEN NEW.status='WITHDRAWN' THEN 'CLIENT_CANCELED'
        WHEN NEW.status='CANCELED' THEN 'PROVIDER_CANCELED'
        ELSE NULL
    END;

    v_commitment_period := CASE
        WHEN NEW.total_hours IS NOT NULL THEN 'HOURLY'
        WHEN NEW.weekly_limit IS NOT NULL THEN 'WEEKLY'
        ELSE NULL
    END;

    v_commitment_period_count := CASE
        WHEN NEW.total_hours IS NOT NULL THEN NEW.total_hours
        WHEN NEW.weekly_limit IS NOT NULL THEN NEW.weekly_limit
        ELSE NULL
    END;

    IF project.id IS NULL OR contract_status IS NULL THEN
        RETURN NEW; -- Exit the function
    END IF;
    
	
    IF contract.id IS NOT NULL THEN
        UPDATE contracts
        SET
            offer_id=NEW.id,
            name=NEW.offer_message,
            description=NEW.offer_message,
            status=COALESCE(contract_status::contract_status, status),
            type=project.payment_type::payment_type::text::contract_type,
            currency=NEW.currency,
            total_amount=COALESCE(NEW.assignment_total, 999),
            payment_type=NEW.payment_mode,
            project_id=NEW.project_id,
            client_id=NEW.recipient_id,
            provider_id=NEW.offerer_id,
            applicant_id=NEW.applicant_id,
            currency_rate=NULL,
            commitment=NULL,
            commitment_period=v_commitment_period::text::contract_commitment_period,
            commitment_period_count=v_commitment_period_count
        WHERE id = contract.id;
    ELSE
        INSERT INTO
        contracts
        (
            offer_id,
            name,
            description,
            status,
            type,
            currency,
            total_amount,
            payment_type,
            project_id,
            client_id,
            provider_id,
            applicant_id,
            currency_rate,
            commitment,
            commitment_period,
            commitment_period_count
        )
        VALUES (
            NEW.id,
            NEW.offer_message,
            NEW.offer_message,
            contract_status::contract_status,
            project.payment_type::payment_type::text::contract_type,
            NEW.currency,
            COALESCE(NEW.assignment_total, NULL),
            NEW.payment_mode,
            NEW.project_id,
            NEW.recipient_id,
            NEW.offerer_id,
            NEW.applicant_id,
            NULL,
            NULL,
            v_commitment_period::text::contract_commitment_period,
            v_commitment_period_count
        );
    END IF;

    RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION upsert_contract_on_mission() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    contract RECORD;
    offer RECORD;
    project RECORD;
    contract_status TEXT;
    v_commitment_period TEXT;
    v_commitment_period_count INTEGER;
BEGIN
    SELECT * INTO contract
    FROM contracts
    WHERE offer_id = NEW.offer_id;

    SELECT * INTO project
    FROM projects
    WHERE id = NEW.project_id AND payment_type IS NOT NULL;

    SELECT * INTO offer
    FROM offers
    WHERE id = NEW.offer_id;

    contract_status := CASE
        WHEN NEW.status='ACTIVE' THEN 'SIGNED'
        WHEN NEW.status='COMPLETE' OR NEW.status='CONFIRMED' THEN 'COMPLETED' --We can't differentiate between these
        WHEN NEW.status='CANCELED' THEN 'CLIENT_CANCELED'
        WHEN NEW.status='KICKED_OUT' THEN 'PROVIDER_CANCELED'
        ELSE NULL
    END;

    
    IF project.id IS NULL OR contract_status IS NULL OR offer.id IS NULL THEN
        RETURN NEW; -- Exit the function
    END IF;

    v_commitment_period := CASE
        WHEN offer.total_hours IS NOT NULL THEN 'HOURLY'
        WHEN offer.weekly_limit IS NOT NULL THEN 'WEEKLY'
        ELSE NULL
    END;

    v_commitment_period_count := CASE
        WHEN offer.total_hours IS NOT NULL THEN offer.total_hours
        WHEN offer.weekly_limit IS NOT NULL THEN offer.weekly_limit
        ELSE NULL
    END;

    IF contract.id IS NOT NULL THEN
        UPDATE contracts
        SET
            mission_id=NEW.id,
            status=contract_status::contract_status
        WHERE id = contract.id;
    ELSE
        INSERT INTO
        contracts
        (
            offer_id,
            mission_id,
            name,
            description,
            status,
            type,
            currency,
            total_amount,
            payment_type,
            project_id,
            client_id,
            provider_id,
            applicant_id,
            currency_rate,
            commitment,
            commitment_period,
            commitment_period_count
        )
        VALUES (
            offer.id,
            NEW.id,
            offer.offer_message,
            offer.offer_message,
            contract_status::contract_status,
            project.payment_type::payment_type::text::contract_type,
            offer.currency,
            COALESCE(offer.assignment_total, NULL),
            offer.payment_mode,
            offer.project_id,
            offer.recipient_id,
            offer.offerer_id,
            offer.applicant_id,
            NULL,
            NULL,
            v_commitment_period::text::contract_commitment_period,
            v_commitment_period_count
        );
    END IF;
    RETURN NEW;
END;
$$;

CREATE OR REPLACE TRIGGER upsert_contract_on_offer AFTER INSERT OR UPDATE ON offers FOR EACH ROW EXECUTE FUNCTION upsert_contract_on_offer();
CREATE OR REPLACE TRIGGER upsert_contract_on_mission AFTER INSERT OR UPDATE ON missions FOR EACH ROW EXECUTE FUNCTION upsert_contract_on_mission();

-- Migrating the Offers and Missions
UPDATE offers SET id=id;
UPDATE missions SET id=id;