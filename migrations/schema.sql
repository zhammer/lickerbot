--
-- PostgreSQL database dump
--

-- Dumped from database version 12.2 (Debian 12.2-1.pgdg100+1)
-- Dumped by pg_dump version 12.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: bootlickers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.bootlickers (
    id uuid NOT NULL,
    twitter_user_id bigint NOT NULL,
    twitter_handle character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.bootlickers OWNER TO postgres;

--
-- Name: licks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.licks (
    id uuid NOT NULL,
    tweet_id bigint NOT NULL,
    tweet_text text NOT NULL,
    bootlicker_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.licks OWNER TO postgres;

--
-- Name: pledged_donations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pledged_donations (
    id uuid NOT NULL,
    bootlicker_id uuid NOT NULL,
    amount integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    CONSTRAINT pledged_donations_amount_check CHECK ((amount > 0))
);


ALTER TABLE public.pledged_donations OWNER TO postgres;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO postgres;

--
-- Name: bootlickers bootlickers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bootlickers
    ADD CONSTRAINT bootlickers_pkey PRIMARY KEY (id);


--
-- Name: bootlickers bootlickers_twitter_handle_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bootlickers
    ADD CONSTRAINT bootlickers_twitter_handle_key UNIQUE (twitter_handle);


--
-- Name: bootlickers bootlickers_twitter_user_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bootlickers
    ADD CONSTRAINT bootlickers_twitter_user_id_key UNIQUE (twitter_user_id);


--
-- Name: licks licks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.licks
    ADD CONSTRAINT licks_pkey PRIMARY KEY (id);


--
-- Name: licks licks_tweet_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.licks
    ADD CONSTRAINT licks_tweet_id_key UNIQUE (tweet_id);


--
-- Name: pledged_donations pledged_donations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pledged_donations
    ADD CONSTRAINT pledged_donations_pkey PRIMARY KEY (id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: licks licks_bootlicker_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.licks
    ADD CONSTRAINT licks_bootlicker_id_fkey FOREIGN KEY (bootlicker_id) REFERENCES public.bootlickers(id) ON DELETE CASCADE;


--
-- Name: pledged_donations pledged_donations_bootlicker_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pledged_donations
    ADD CONSTRAINT pledged_donations_bootlicker_id_fkey FOREIGN KEY (bootlicker_id) REFERENCES public.bootlickers(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

