-- Add Offer and Mission to contract
ALTER TABLE contracts
    ADD COLUMN offer_id UUID,
    ADD COLUMN mission_id UUID,
    ADD CONSTRAINT fk_offer FOREIGN KEY (offer_id) REFERENCES offers(id) ON DELETE SET NULL,
    ADD CONSTRAINT fk_mission FOREIGN KEY (mission_id) REFERENCES missions(id) ON DELETE SET NULL;

-- Add triggers
CREATE OR REPLACE FUNCTION upsert_contract_on_offer() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    contract RECORD;
    project RECORD;
    contract_status TEXT;
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
	
    IF contract.id IS NOT NULL AND project.id IS NOT NULL THEN
        UPDATE contracts
        SET
            offer_id=NEW.id,
            name=NEW.offer_message,-- Is it same as description?
            description=NEW.offer_message,
            status=COALESCE(contract_status::contract_status, status),
            type=project.payment_type::payment_type::text::contract_type,
            currency=NEW.currency,
            total_amount=NEW.assignment_total,
            payment_type=NEW.payment_mode,
            project_id=NEW.project_id,
            client_id=NEW.recipient_id,
            provider_id=NEW.offerer_id,
            applicant_id=NEW.applicant_id,
            commitment_period='HOURLY' -- Couldn't find the corresponding column?
        WHERE id = contract.id;
    ELSEIF project.id IS NOT NULL THEN
        INSERT INTO
        contracts
        (
            offer_id,
            name,--?
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
            commitment_period
        )
        VALUES (
            NEW.id,
            NEW.offer_message,--?
            NEW.offer_message,
            contract_status::contract_status,
            project.payment_type::payment_type::text::contract_type,
            NEW.currency,
            NEW.assignment_total,
            NEW.payment_mode,
            NEW.project_id,
            NEW.recipient_id,
            NEW.offerer_id,
            NEW.applicant_id,
            'HOURLY' --?
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
        WHEN NEW.status='ACTIVE' OR NEW.status='COMPLETE' OR NEW.status='CONFIRMED' THEN 'SIGNED'
        WHEN NEW.status='CANCELED' THEN 'CLIENT_CANCELED'
        WHEN NEW.status='KICKED_OUT' THEN 'PROVIDER_CANCELED'
        ELSE NULL
    END;

    IF contract.id IS NOT NULL THEN
        UPDATE contracts
        SET
            mission_id=NEW.id,
            status=contract_status::contract_status
        WHERE id = contract.id;
    ELSEIF offer.id IS NOT NULL AND project.id IS NOT NULL THEN
        INSERT INTO
        contracts
        (
            offer_id,
            mission_id,
            name, --?
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
            commitment_period --?
        )
        VALUES (
            offer.id,
            NEW.id,
            offer.offer_message, --?
            offer.offer_message,
            contract_status::contract_status,
            project.payment_type::payment_type::text::contract_type,
            offer.currency,
            offer.assignment_total,
            offer.payment_mode,
            offer.project_id,
            offer.recipient_id,
            offer.offerer_id,
            offer.applicant_id,
            'HOURLY' --?
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