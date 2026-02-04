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
        '早起きマスター',
        '5日連続で10時前にチェックイン',
        'time',
        'early_check_in',
        '{"days": 5, "time": "10:00:00"}',
        10,
        1
    ),
    (
        'early_bird_30',
        '早起きの人',
        '10日連続で10時前にチェックイン',
        'time',
        'early_check_in',
        '{"days": 10, "time": "10:00:00"}',
        50,
        2
    ),
    -- 夜型系
    (
        'night_owl_10',
        '夜に来るのもいいよね',
        '20時以降に5回チェックイン',
        'time',
        'late_check_in',
        '{"count": 5, "time": "20:00:00"}',
        10,
        10
    ),
    (
        'night_owl_50',
        '生活習慣大丈夫？',
        '22時以降に5回チェックイン',
        'time',
        'late_check_in',
        '{"count": 10, "time": "22:00:00"}',
        10,
        11
    ),
    -- ストリーク系
    (
        'streak_7',
        '7日連続で来るなんて...',
        '7日連続で研究室に来る',
        'streak',
        'streak_days',
        '{"days": 7}',
        20,
        20
    ),
    (
        'streak_14',
        '2週間連続で来るなんて...',
        '14日連続で研究室に来る',
        'streak',
        'streak_days',
        '{"days": 14}',
        40,
        21
    ),
    (
        'streak_30',
        '勤勉マン',
        '30日連続で研究室に来る',
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
        '30日きたよ',
        '30日間研究室に来た',
        'attendance',
        'total_days',
        '{"days": 30}',
        20,
        40
    ),
    (
        'attendance_60',
        '60日きたよ',
        '60日間研究室に来た',
        'attendance',
        'total_days',
        '{"days": 60}',
        50,
        41
    ),
    (
        'attendance_100',
        '100日きたよ',
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
        '土日に来るのもいいよね',
        '土日に5回チェックイン',
        'special',
        'weekend_check_in',
        '{"count": 5}',
        10,
        101
    );