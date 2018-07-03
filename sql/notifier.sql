--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.5
-- Dumped by pg_dump version 9.6.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: notifications; Type: TABLE; Schema: public; Owner: notifier
--

CREATE TABLE notifications (
    nid integer NOT NULL,
    type character varying(256) DEFAULT ''::character varying NOT NULL,
    uid integer NOT NULL,
    timestamp_add timestamp with time zone DEFAULT now() NOT NULL,
    value character varying(256) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE notifications OWNER TO notifier;

--
-- Name: notifications_nid_seq; Type: SEQUENCE; Schema: public; Owner: notifier
--

CREATE SEQUENCE notifications_nid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE notifications_nid_seq OWNER TO notifier;

--
-- Name: notifications_nid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: notifier
--

ALTER SEQUENCE notifications_nid_seq OWNED BY notifications.nid;


--
-- Name: notifications nid; Type: DEFAULT; Schema: public; Owner: notifier
--

ALTER TABLE ONLY notifications ALTER COLUMN nid SET DEFAULT nextval('notifications_nid_seq'::regclass);


--
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: notifier
--

COPY notifications (nid, type, uid, timestamp_add, value) FROM stdin;
\.


--
-- Name: notifications_nid_seq; Type: SEQUENCE SET; Schema: public; Owner: notifier
--

SELECT pg_catalog.setval('notifications_nid_seq', 1, false);


--
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: notifier
--

ALTER TABLE ONLY notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (nid);


--
-- Name: public; Type: ACL; Schema: -; Owner: notifier
--

REVOKE ALL ON SCHEMA public FROM postgres;
REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO notifier;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

