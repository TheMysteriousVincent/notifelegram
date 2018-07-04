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
-- Name: commits; Type: TABLE; Schema: public; Owner: notifier
--

CREATE TABLE commits (
    commitid integer NOT NULL,
    timestamp_add timestamp with time zone DEFAULT now() NOT NULL,
    "chatId" bigint DEFAULT 0 NOT NULL
);


ALTER TABLE commits OWNER TO notifier;

--
-- Name: commits_commitid_seq; Type: SEQUENCE; Schema: public; Owner: notifier
--

CREATE SEQUENCE commits_commitid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE commits_commitid_seq OWNER TO notifier;

--
-- Name: commits_commitid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: notifier
--

ALTER SEQUENCE commits_commitid_seq OWNED BY commits.commitid;


--
-- Name: mentions; Type: TABLE; Schema: public; Owner: notifier
--

CREATE TABLE mentions (
    mentionid integer NOT NULL,
    "chatId" bigint DEFAULT 0 NOT NULL,
    timestamp_add timestamp with time zone DEFAULT now() NOT NULL,
    gitlabusername character varying(256) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE mentions OWNER TO notifier;

--
-- Name: mentions_mentionid_seq; Type: SEQUENCE; Schema: public; Owner: notifier
--

CREATE SEQUENCE mentions_mentionid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE mentions_mentionid_seq OWNER TO notifier;

--
-- Name: mentions_mentionid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: notifier
--

ALTER SEQUENCE mentions_mentionid_seq OWNED BY mentions.mentionid;


--
-- Name: commits commitid; Type: DEFAULT; Schema: public; Owner: notifier
--

ALTER TABLE ONLY commits ALTER COLUMN commitid SET DEFAULT nextval('commits_commitid_seq'::regclass);


--
-- Name: mentions mentionid; Type: DEFAULT; Schema: public; Owner: notifier
--

ALTER TABLE ONLY mentions ALTER COLUMN mentionid SET DEFAULT nextval('mentions_mentionid_seq'::regclass);


--
-- Data for Name: commits; Type: TABLE DATA; Schema: public; Owner: notifier
--

COPY commits (commitid, timestamp_add, "chatId") FROM stdin;
1	2018-07-04 11:00:06.309049+02	406907138
\.


--
-- Name: commits_commitid_seq; Type: SEQUENCE SET; Schema: public; Owner: notifier
--

SELECT pg_catalog.setval('commits_commitid_seq', 1, true);


--
-- Data for Name: mentions; Type: TABLE DATA; Schema: public; Owner: notifier
--

COPY mentions (mentionid, "chatId", timestamp_add, gitlabusername) FROM stdin;
\.


--
-- Name: mentions_mentionid_seq; Type: SEQUENCE SET; Schema: public; Owner: notifier
--

SELECT pg_catalog.setval('mentions_mentionid_seq', 3, true);


--
-- Name: commits commits_pkey; Type: CONSTRAINT; Schema: public; Owner: notifier
--

ALTER TABLE ONLY commits
    ADD CONSTRAINT commits_pkey PRIMARY KEY (commitid);


--
-- Name: mentions mentions_pkey; Type: CONSTRAINT; Schema: public; Owner: notifier
--

ALTER TABLE ONLY mentions
    ADD CONSTRAINT mentions_pkey PRIMARY KEY (mentionid);


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

