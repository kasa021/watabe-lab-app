-- 称号の初期データ
INSERT INTO achievements (
        code,
        name,
        description,
        category,
        condition_type,
        condition_value,
        points_reward,
        display_order
    )
VALUES -- 早起き系
    (
        'early_bird_7',
        '早起き博士',
        '7日連続で9時前にチェックイン',
        'time',
        'early_check_in',
        '{"days": 7, "time": "09:00:00"}',
        10,
        1
    ),
    (
        'early_bird_30',
        '早起きマスター',
        '30日連続で9時前にチェックイン',
        'time',
        'early_check_in',
        '{"days": 30, "time": "09:00:00"}',
        50,
        2
    ),
    -- 夜型系
    (
        'night_owl_10',
        '夜型研究者',
        '22時以降に10回チェックイン',
        'time',
        'late_check_in',
        '{"count": 10, "time": "22:00:00"}',
        10,
        10
    ),
    (
        'night_owl_50',
        '深夜の住人',
        '22時以降に50回チェックイン',
        'time',
        'late_check_in',
        '{"count": 50, "time": "22:00:00"}',
        50,
        11
    ),
    -- ストリーク系
    (
        'streak_7',
        '7日連続出席',
        '7日連続で研究室に来る',
        'streak',
        'streak_days',
        '{"days": 7}',
        20,
        20
    ),
    (
        'streak_14',
        '2週間連続出席',
        '14日連続出席',
        'streak',
        'streak_days',
        '{"days": 14}',
        40,
        21
    ),
    (
        'streak_30',
        '皆勤賞',
        '30日連続出席',
        'streak',
        'streak_days',
        '{"days": 30}',
        100,
        22
    ),
    -- 長時間滞在系
    (
        'total_50h',
        '50時間達成',
        '累計滞在時間50時間',
        'time',
        'total_hours',
        '{"hours": 50}',
        10,
        30
    ),
    (
        'total_100h',
        '100時間達成',
        '累計滞在時間100時間',
        'time',
        'total_hours',
        '{"hours": 100}',
        30,
        31
    ),
    (
        'total_300h',
        '300時間達成',
        '累計滞在時間300時間',
        'time',
        'total_hours',
        '{"hours": 300}',
        100,
        32
    ),
    (
        'total_500h',
        '研究室の主',
        '累計滞在時間500時間',
        'time',
        'total_hours',
        '{"hours": 500}',
        200,
        33
    ),
    -- 出席日数系
    (
        'attendance_30',
        '30日出席',
        '30日間研究室に来た',
        'attendance',
        'total_days',
        '{"days": 30}',
        20,
        40
    ),
    (
        'attendance_60',
        '60日出席',
        '60日間研究室に来た',
        'attendance',
        'total_days',
        '{"days": 60}',
        50,
        41
    ),
    (
        'attendance_100',
        '100日出席',
        '100日間研究室に来た',
        'attendance',
        'total_days',
        '{"days": 100}',
        100,
        42
    ),
    -- 特別系
    (
        'first_check_in',
        'はじめの一歩',
        '初めてチェックイン',
        'special',
        'first_time',
        '{}',
        5,
        100
    ),
    (
        'weekend_warrior',
        '週末戦士',
        '土日に10回チェックイン',
        'special',
        'weekend_check_in',
        '{"count": 10}',
        30,
        101
    );