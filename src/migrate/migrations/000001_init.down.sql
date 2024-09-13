BEGIN;

DROP TABLE IF EXISTS bid_decision;

DROP TABLE IF EXISTS bid_review;

DROP TABLE IF EXISTS bid_information;

DROP TABLE IF EXISTS bid;

DROP TABLE IF EXISTS tender_information;

DROP TABLE IF EXISTS tender;

DROP TYPE IF EXISTS bid_decision_type;

DROP TYPE IF EXISTS bid_author_type;

DROP TYPE IF EXISTS bid_status;

DROP TYPE IF EXISTS tender_service_type;

DROP TYPE IF EXISTS tender_status;

COMMIT;
