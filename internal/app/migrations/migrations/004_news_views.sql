-- +goose Up

ALTER TABLE IF EXISTS public.news
    ADD COLUMN IF NOT EXISTS views integer DEFAULT 0;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION get_news_with_views(news_id integer)
    RETURNS TABLE (
                      header character(80),
                      description character(250),
                      body text,
                      date date,
                      preview_url character(256),
                      author_name character(100),
                      ai_correct boolean,
                      views integer
                  )
AS $$
BEGIN
    -- Увеличиваем значение поля "views" на 1
    EXECUTE 'UPDATE public.news
           SET views = views + 1
           WHERE id = ' || news_id;

    -- Выполняем SELECT запрос и возвращаем результаты
    RETURN QUERY
        SELECT n.header, n.description, n.body, n.date, n.preview_url, a.name, n.ai_correct, n.views
        FROM public.news AS n
                 LEFT JOIN public.news_authors AS a ON n.author = a.id
        WHERE n.id = news_id;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down

ALTER TABLE IF EXISTS public.news
    DROP COLUMN views;