CREATE OR REPLACE FUNCTION increment_work_count()
RETURNS TRIGGER AS $$
BEGIN
    -- 开始一个事务块
    BEGIN
        UPDATE users
        SET work_count = work_count + 1
        WHERE id = NEW.author;

        -- 如果更新失败，抛出异常
        IF NOT FOUND THEN
            RAISE EXCEPTION '更新work_count计数失败';
        END IF;
    EXCEPTION
        WHEN others THEN
            -- 出现异常时回滚事务
            RAISE NOTICE '更新work_count计数时发生异常: %', SQLERRM;
            ROLLBACK;
    END;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION decrement_work_count()
RETURNS TRIGGER AS $$
BEGIN
    -- 开始一个事务块
    BEGIN
        UPDATE users
        SET work_count = work_count - 1
        WHERE id = OLD.author;

        -- 如果更新失败，抛出异常
        IF NOT FOUND THEN
            RAISE EXCEPTION '更新work_count计数失败';
        END IF;
    EXCEPTION
        WHEN others THEN
            -- 出现异常时回滚事务
            RAISE NOTICE '更新work_count计数时发生异常: %', SQLERRM;
            ROLLBACK;
    END;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER increment_work_count_trigger
-- 插入后触发
AFTER INSERT ON videos
FOR EACH ROW
EXECUTE FUNCTION increment_work_count();

CREATE OR REPLACE TRIGGER decrement_work_count_trigger
-- 删除后触发
AFTER DELETE ON videos
FOR EACH ROW
EXECUTE FUNCTION decrement_work_count();
