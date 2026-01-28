-- システム設定の初期データ
INSERT INTO settings (key, value, description)
VALUES (
        'auto_checkout_minutes',
        '120',
        '自動チェックアウトまでの時間（分）'
    ),
    (
        'holidays',
        '["2025-01-01", "2025-01-13", "2025-02-11", "2025-02-23", "2025-03-20", "2025-04-29", "2025-05-03", "2025-05-04", "2025-05-05", "2025-07-21", "2025-08-11", "2025-09-15", "2025-09-23", "2025-10-13", "2025-11-03", "2025-11-23", "2025-11-24"]',
        '休日リスト（2025年の祝日）'
    ),
    (
        'allowed_ip_range',
        '{"ips": ["133.38.201.125", "172.19.0.1"]}',
        '許可されたIP範囲'
    ),
    (
        'gps_location',
        '{"latitude": 35.862934, "longitude": 139.607886, "radius_meters": 200}',
        '研究室の位置情報と許容範囲'
    )
ON CONFLICT (key) DO UPDATE 
SET value = EXCLUDED.value, description = EXCLUDED.description;