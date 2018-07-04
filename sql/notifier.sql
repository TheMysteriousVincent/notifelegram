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
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: commits; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE commits (
    commitid integer NOT NULL,
    timestamp_add timestamp with time zone DEFAULT now() NOT NULL,
    "chatId" bigint DEFAULT 0 NOT NULL
);


--
-- Name: commits_commitid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE commits_commitid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: commits_commitid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE commits_commitid_seq OWNED BY commits.commitid;


--
-- Name: mentions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE mentions (
    mentionid integer NOT NULL,
    "chatId" bigint DEFAULT 0 NOT NULL,
    timestamp_add timestamp with time zone DEFAULT now() NOT NULL,
    gitlabusername character varying(256) DEFAULT ''::character varying NOT NULL
);


--
-- Name: mentions_mentionid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE mentions_mentionid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: mentions_mentionid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE mentions_mentionid_seq OWNED BY mentions.mentionid;


--
-- Name: commits commitid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY commits ALTER COLUMN commitid SET DEFAULT nextval('commits_commitid_seq'::regclass);


--
-- Name: mentions mentionid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY mentions ALTER COLUMN mentionid SET DEFAULT nextval('mentions_mentionid_seq'::regclass);


--
-- Name: commits commits_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY commits
    ADD CONSTRAINT commits_pkey PRIMARY KEY (commitid);


--
-- Name: mentions mentions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY mentions
    ADD CONSTRAINT mentions_pkey PRIMARY KEY (mentionid);


--
-- PostgreSQL database dump complete
--

