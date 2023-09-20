CREATE OR REPLACE FUNCTION increment_follow_count()
RETURNS TRIGGER AS $$
BEGIN
    -- 开始一个事务块
    BEGIN
        UPDATE users
        SET follow_count = follow_count + 1
        WHERE id = NEW.user_id;

        UPDATE users
        SET follower_count = follower_count + 1
        WHERE id = NEW.follow_id;

        -- 如果更新失败，抛出异常
        IF NOT FOUND THEN
            RAISE EXCEPTION '更新follow_count计数失败';
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

CREATE OR REPLACE FUNCTION decrement_follow_count()
RETURNS TRIGGER AS $$
BEGIN
    -- 开始一个事务块
    BEGIN
        UPDATE users
        SET follower_count = follower_count - 1
        WHERE id = OLD.user_id;

        UPDATE users
        SET follower_count = follower_count - 1
        WHERE id = OLD.follow_id;

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

CREATE OR REPLACE TRIGGER increment_follow_count_trigger
-- 插入后触发
AFTER INSERT ON follows
FOR EACH ROW
EXECUTE FUNCTION increment_follow_count();

CREATE OR REPLACE TRIGGER decrement_follow_count_trigger
-- 删除后触发
AFTER DELETE ON follows
FOR EACH ROW
EXECUTE FUNCTION decrement_follow_count();
