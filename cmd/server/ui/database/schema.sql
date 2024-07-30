--++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
--1) create tables
--++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

DROP TABLE IF EXISTS posts, authors;

CREATE TABLE authors (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    author_id BIGINT REFERENCES authors(id) NOT NULL,
    title TEXT  NOT NULL,
    content TEXT NOT NULL,
    created_at BIGINT NOT NULL,
    published_at BIGINT NOT NULL
);




--++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
--2) create functions
--++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

DROP FUNCTION IF EXISTS authors_func_delete, authors_func_update, authors_func_insert, authors_func_view;
DROP FUNCTION IF EXISTS posts_func_delete, posts_func_update, posts_func_insert, posts_func_view;

--=======================
--table: authors
--=======================
--insert
CREATE FUNCTION authors_func_insert(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	INSERT INTO authors (name) VALUES ((json_data ->> 'name')::TEXT) RETURNING id INTO new_id;
	
	IF new_id IS NULL THEN
		RAISE EXCEPTION 'Parameter value cannot be null';
	END IF;

	SELECT json_build_object('id',new_id,'err','') INTO json_result;
  	RETURN json_result;

EXCEPTION
    WHEN others THEN
		GET STACKED DIAGNOSTICS err_context = PG_EXCEPTION_CONTEXT;
    	GET STACKED DIAGNOSTICS err_mess = MESSAGE_TEXT;

        SELECT json_build_object('id',null,'err',err_mess||err_context) INTO json_result;
  		RETURN json_result;  
END;
$$ LANGUAGE plpgsql;

--update
CREATE FUNCTION authors_func_update(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	UPDATE authors SET name = (json_data ->> 'name')::TEXT WHERE id = (json_data ->> 'id')::BIGINT RETURNING id INTO new_id; 
	
	IF new_id IS NULL THEN
		RAISE EXCEPTION 'Parameter value cannot be null';
	END IF;

	SELECT json_build_object('id',new_id,'err','') INTO json_result;
  	RETURN json_result;

EXCEPTION
    WHEN others THEN
		GET STACKED DIAGNOSTICS err_context = PG_EXCEPTION_CONTEXT;
    	GET STACKED DIAGNOSTICS err_mess = MESSAGE_TEXT;

        SELECT json_build_object('id',null,'err',err_mess||err_context) INTO json_result;
  		RETURN json_result;  
END;
$$ LANGUAGE plpgsql;


--delete
CREATE FUNCTION authors_func_delete(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	DELETE FROM authors WHERE id = (json_data ->> 'id')::BIGINT; 

	SELECT json_build_object('id',(json_data ->> 'id')::BIGINT,'err','') INTO json_result;
  	RETURN json_result;

EXCEPTION
    WHEN others THEN
		GET STACKED DIAGNOSTICS err_context = PG_EXCEPTION_CONTEXT;
    	GET STACKED DIAGNOSTICS err_mess = MESSAGE_TEXT;

        SELECT json_build_object('id',null,'err',err_mess||err_context) INTO json_result;
  		RETURN json_result;  
END;
$$ LANGUAGE plpgsql;

--select view
CREATE FUNCTION authors_func_view(
		json_data jsonb
) 
RETURNS TABLE (
	id BIGINT,
    name TEXT
) AS $$
DECLARE
  	par_id BIGINT = 0;
    par_name TEXT = '';
BEGIN

	IF (json_data ->> 'id') IS NOT NULL THEN
		par_id = (json_data ->> 'id')::BIGINT;
	END IF;

	IF (json_data ->> 'name') IS NOT NULL THEN
		par_name = (json_data ->> 'name')::TEXT;
	END IF;


	RETURN QUERY
		SELECT authors.id,
			   authors.name   
		FROM authors
		WHERE
			(par_id = 0 OR authors.id = par_id) AND
			(par_name = '' OR authors.name LIKE '%'||par_name||'%')
		ORDER BY authors.id;
	
END;
$$ LANGUAGE plpgsql;



--=======================
--table: posts
--=======================
--insert
CREATE FUNCTION posts_func_insert(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	INSERT INTO posts (
		author_id, 
		title, 
		content,
		created_at,
		published_at
		) 
	VALUES (
		(json_data ->> 'author_id')::BIGINT, 
		(json_data ->> 'title')::TEXT, 
		(json_data ->> 'content')::TEXT, 
		(json_data ->> 'created_at')::BIGINT,
		(json_data ->> 'published_at')::BIGINT
		)
	RETURNING id INTO new_id;
	
	IF new_id IS NULL THEN
		RAISE EXCEPTION 'Parameter value cannot be null. ';
	END IF;

	SELECT json_build_object('id',new_id,'err','') INTO json_result;
  	RETURN json_result;

EXCEPTION
    WHEN others THEN
		GET STACKED DIAGNOSTICS err_context = PG_EXCEPTION_CONTEXT;
    	GET STACKED DIAGNOSTICS err_mess = MESSAGE_TEXT;

        SELECT json_build_object('id',null,'err',err_mess||err_context) INTO json_result;
  		RETURN json_result;  
END;
$$ LANGUAGE plpgsql;

--update
CREATE FUNCTION posts_func_update(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
  	new_id BIGINT;
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN
    
	UPDATE posts SET 
		author_id = CASE WHEN (json_data ->> 'author_id') IS NOT NULL THEN (json_data ->> 'author_id')::BIGINT ELSE author_id END,
		title = CASE WHEN (json_data ->> 'title') IS NOT NULL THEN (json_data ->> 'title')::TEXT ELSE title END,
		content = CASE WHEN (json_data ->> 'content') IS NOT NULL THEN (json_data ->> 'content')::TEXT ELSE content END,
		created_at = CASE WHEN (json_data ->> 'created_at') IS NOT NULL THEN (json_data ->> 'created_at')::BIGINT ELSE created_at END,
		published_at = CASE WHEN (json_data ->> 'published_at') IS NOT NULL THEN (json_data ->> 'published_at')::BIGINT ELSE published_at END
	WHERE 
		id = (json_data ->> 'id')::BIGINT
	RETURNING id INTO new_id;
	
	
	IF new_id IS NULL THEN
		RAISE EXCEPTION 'Parameter value cannot be null. A closed task cannot be modified. ';
	END IF;

	SELECT json_build_object('id',new_id,'err','') INTO json_result;
  	RETURN json_result;

EXCEPTION
    WHEN others THEN
		GET STACKED DIAGNOSTICS err_context = PG_EXCEPTION_CONTEXT;
    	GET STACKED DIAGNOSTICS err_mess = MESSAGE_TEXT;

        SELECT json_build_object('id',null,'err',err_mess||err_context) INTO json_result;
  		RETURN json_result;  
END;
$$ LANGUAGE plpgsql;


--delete
CREATE FUNCTION posts_func_delete(
		json_data jsonb
) 
RETURNS jsonb AS $$
DECLARE
	err_mess TEXT;
	err_context TEXT;
	json_result jsonb;
BEGIN

	DELETE FROM posts WHERE id = (json_data ->> 'id')::BIGINT; 

	SELECT json_build_object('id',(json_data ->> 'id')::BIGINT,'err','') INTO json_result;
  	RETURN json_result;

EXCEPTION
    WHEN others THEN
		GET STACKED DIAGNOSTICS err_context = PG_EXCEPTION_CONTEXT;
    	GET STACKED DIAGNOSTICS err_mess = MESSAGE_TEXT;

        SELECT json_build_object('id',null,'err',err_mess||err_context) INTO json_result;
  		RETURN json_result;  
END;
$$ LANGUAGE plpgsql;

--select
CREATE FUNCTION posts_func_view(
		json_data jsonb
) 
RETURNS TABLE (
	id BIGINT,
    author_id BIGINT, 
	author_name TEXT, 
    title TEXT, 
	content TEXT, 
    created_at BIGINT, 
	created_at_txt TEXT, 
    published_at BIGINT, 
    published_at_txt TEXT
) AS $$
DECLARE
  	par_id BIGINT = 0;
    par_author_id BIGINT = 0;
    par_title TEXT = '';
    par_content TEXT = '';
	par_created_at BIGINT = 0;
    par_published_at BIGINT = 0;
BEGIN

	IF (json_data ->> 'id') IS NOT NULL THEN
		par_id = (json_data ->> 'id')::BIGINT;
	END IF;

	IF (json_data ->> 'author_id') IS NOT NULL THEN
		par_author_id = (json_data ->> 'author_id')::BIGINT;
	END IF;

	IF (json_data ->> 'title') IS NOT NULL THEN
		par_title = (json_data ->> 'title')::TEXT;
	END IF;

	IF (json_data ->> 'content') IS NOT NULL THEN
		par_content = (json_data ->> 'content')::TEXT;
	END IF;

	IF (json_data ->> 'created_at') IS NOT NULL THEN
		par_created_at = (json_data ->> 'created_at')::BIGINT;
	END IF;

	IF (json_data ->> 'published_at') IS NOT NULL THEN
		par_published_at = (json_data ->> 'published_at')::BIGINT;
	END IF;

	
	RETURN QUERY
		SELECT 
			posts.id as id,
			posts.author_id as author_id,
			COALESCE((
				SELECT authors.name 
				FROM   authors 
				WHERE  authors.id = posts.author_id
			), '') as author_name,
			posts.title as title,
			posts.content as content,
			posts.created_at as created_at,
			COALESCE(TO_CHAR(TO_TIMESTAMP(posts.created_at/1000), 'DD.MM.YYYY HH24:MI:SS'), '') as created_at_txt, 
			posts.published_at as published_at,
			COALESCE(TO_CHAR(TO_TIMESTAMP(posts.published_at/1000), 'DD.MM.YYYY HH24:MI:SS'), '') as published_at_txt
		FROM 
			posts
		WHERE
			(par_id = 0 OR posts.id = par_id) AND
			(par_author_id = 0 OR posts.author_id = par_author_id) AND 
			(par_title = '' OR posts.title LIKE '%'||par_title||'%') AND
			(par_content = '' OR posts.content LIKE '%'||par_content||'%') AND
			(par_created_at = 0 OR to_timestamp(posts.created_at / 1000)::date = to_timestamp(par_created_at / 1000)::date) AND   
			(par_published_at = 0 OR to_timestamp(posts.published_at / 1000)::date = to_timestamp(par_published_at / 1000)::date)
		ORDER BY 
			posts.id;
END;
$$ LANGUAGE plpgsql;


--++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
--3) create test data
--++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

/*
DO $$
DECLARE
  	i INT;
	str001 TEXT;
	str002 TEXT;
	bi001 TEXT;
	bi002 TEXT;
BEGIN
	
	--authors: insert
	i:=1;
  	WHILE i <= 100 LOOP
		str001 := 'Authors_' || LPAD(i::TEXT, 10, '0');
    	PERFORM * FROM authors_func_insert(('{
				  "name": "'||str001||'"
				}')::jsonb);
    	i := i + 1;
  	END LOOP;

	--posts: insert
	i:=1;
	WHILE i <= 50 LOOP
		str001 := 'Title_' || LPAD(i::TEXT, 10, '0');
		str002 := 'Content_' || LPAD(i::TEXT, 10, '0');
		bi001 := (extract(epoch from now())::BIGINT- (i*2 * 24 * 60 * 60))*1000;
		bi002 := (extract(epoch from now())::BIGINT- (i*1 * 24 * 60 * 60))*1000;
    	PERFORM * FROM posts_func_insert(('{
		  "author_id": '||i||',
		  "title": "'||str001||'",
		  "content": "'||str002||'",
		  "created_at": '||bi001||',
		  "published_at": '||bi002||'
		}')::jsonb);
    	i := i + 1;
  	END LOOP;


END;
$$;*/

--select * from authors_func_view(('{"id": 0, "name": "01"}')::jsonb);
--select * from posts_func_view(('{"id": 0, "content": "01", "published_at": 1720940691000}')::jsonb);
