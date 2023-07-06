-- BEGIN TRANSACTION ;

	
-- 	CREATE TABLE IF NOT EXISTS users
-- 	(
-- 		id bigint NOT NULL,
-- 		username text COLLATE pg_catalog."default" NOT NULL,
-- 		password text COLLATE pg_catalog."default" NOT NULL,
-- 		CONSTRAINT urls_pkey PRIMARY KEY (id)
-- 	);
	
-- 	ALTER TABLE public.urls ADD CONSTRAINT original_url_unq UNIQUE(original_url);

-- 	--creating index for text search via short url
-- 	CREATE INDEX urls_short_url_key ON public.urls USING HASH (short_url);
	
-- 	--creating index for text search via original url
-- 	CREATE INDEX urls_original_url_key ON public.urls USING HASH (original_url);


-- COMMIT ;