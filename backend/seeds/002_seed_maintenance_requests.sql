-- ============================================================
-- SEED: 002_seed_maintenance_requests.sql
-- Description: Insert sample maintenance requests with various
--              statuses and priorities for all operator users.
-- ============================================================

INSERT INTO maintenance_requests (
    id, asset_id, problem_description, priority, status,
    created_by, reviewed_by, reviewed_at, review_note,
    created_at, updated_at
)
VALUES
    -- Request 1: Submitted (not reviewed yet) — by operator1
    (
        'a0000000-0000-0000-0000-000000000001',
        'MCH-001',
        'Mesin press hydraulik unit A mengalami kebocoran oli pada selang bagian bawah. Perlu penggantian selang dan pengisian ulang cairan hydraulik.',
        'High',
        'Submitted',
        '00000000-0000-0000-0000-000000000003', -- operator1
        NULL, NULL, NULL,
        NOW() - INTERVAL '2 days',
        NOW() - INTERVAL '2 days'
    ),
    -- Request 2: Approved — by operator1, reviewed by supervisor
    (
        'a0000000-0000-0000-0000-000000000002',
        'CNC-003',
        'Spindle motor CNC mesin 3 mengeluarkan suara tidak normal saat beroperasi di atas 3000 RPM. Kemungkinan bearing aus.',
        'Urgent',
        'Approved',
        '00000000-0000-0000-0000-000000000003', -- operator1
        '00000000-0000-0000-0000-000000000002', -- supervisor
        NOW() - INTERVAL '1 day',
        'Segera jadwalkan teknisi untuk penggantian bearing.',
        NOW() - INTERVAL '3 days',
        NOW() - INTERVAL '1 day'
    ),
    -- Request 3: Rejected — by operator2, reviewed by supervisor
    (
        'a0000000-0000-0000-0000-000000000003',
        'PUMP-007',
        'Pompa air WTP unit 7 terasa kurang kencang saat dipompa.',
        'Low',
        'Rejected',
        '00000000-0000-0000-0000-000000000004', -- operator2
        '00000000-0000-0000-0000-000000000002', -- supervisor
        NOW() - INTERVAL '12 hours',
        'Deskripsi masalah kurang spesifik. Mohon cantumkan tekanan air aktual vs spesifikasi dan sudah berapa lama terjadi.',
        NOW() - INTERVAL '4 days',
        NOW() - INTERVAL '12 hours'
    ),
    -- Request 4: Submitted — by operator2 (with Medium priority)
    (
        'a0000000-0000-0000-0000-000000000004',
        'CONV-012',
        'Belt conveyor line B mengalami slip berulang kali terutama saat membawa beban penuh. Posisi belt sudah bergeser sekitar 2 cm ke kiri.',
        'Medium',
        'Submitted',
        '00000000-0000-0000-0000-000000000004', -- operator2
        NULL, NULL, NULL,
        NOW() - INTERVAL '6 hours',
        NOW() - INTERVAL '6 hours'
    ),
    -- Request 5: Submitted — by operator1 (with Urgent priority)
    (
        'a0000000-0000-0000-0000-000000000005',
        'COMP-002',
        'Kompresor udara unit 2 mati mendadak dan tidak bisa dinyalakan kembali. Panel kontrol menunjukkan error code E-04 (overtemperature). Produksi line C terhenti karena tidak ada supply udara.',
        'Urgent',
        'Submitted',
        '00000000-0000-0000-0000-000000000003', -- operator1
        NULL, NULL, NULL,
        NOW() - INTERVAL '1 hour',
        NOW() - INTERVAL '1 hour'
    ),
    -- Request 6: Approved — by operator2, reviewed by admin
    (
        'a0000000-0000-0000-0000-000000000006',
        'ELEC-015',
        'Panel listrik MCC-3 terdapat MCB yang sering trip tanpa sebab jelas. Suhu panel tinggi sekitar 45°C.',
        'High',
        'Approved',
        '00000000-0000-0000-0000-000000000004', -- operator2
        '00000000-0000-0000-0000-000000000001', -- admin
        NOW() - INTERVAL '2 hours',
        'Disetujui. Tim electrical sudah dijadwalkan untuk inspeksi besok pagi.',
        NOW() - INTERVAL '5 days',
        NOW() - INTERVAL '2 hours'
    )
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- SEED: Audit trail logs for the requests above
-- ============================================================

INSERT INTO request_status_logs (id, request_id, changed_by, from_status, to_status, note, created_at)
VALUES
    -- Log for Request 1: Created as Submitted
    (
        gen_random_uuid(),
        'a0000000-0000-0000-0000-000000000001',
        '00000000-0000-0000-0000-000000000003',
        NULL,
        'Submitted',
        'Request created',
        NOW() - INTERVAL '2 days'
    ),
    -- Log for Request 2: Created, then Approved
    (
        gen_random_uuid(),
        'a0000000-0000-0000-0000-000000000002',
        '00000000-0000-0000-0000-000000000003',
        NULL,
        'Submitted',
        'Request created',
        NOW() - INTERVAL '3 days'
    ),
    (
        gen_random_uuid(),
        'a0000000-0000-0000-0000-000000000002',
        '00000000-0000-0000-0000-000000000002',
        'Submitted',
        'Approved',
        'Segera jadwalkan teknisi untuk penggantian bearing.',
        NOW() - INTERVAL '1 day'
    ),
    -- Log for Request 3: Created, then Rejected
    (
        gen_random_uuid(),
        'a0000000-0000-0000-0000-000000000003',
        '00000000-0000-0000-0000-000000000004',
        NULL,
        'Submitted',
        'Request created',
        NOW() - INTERVAL '4 days'
    ),
    (
        gen_random_uuid(),
        'a0000000-0000-0000-0000-000000000003',
        '00000000-0000-0000-0000-000000000002',
        'Submitted',
        'Rejected',
        'Deskripsi masalah kurang spesifik.',
        NOW() - INTERVAL '12 hours'
    ),
    -- Log for Request 4: Created as Submitted
    (
        gen_random_uuid(),
        'a0000000-0000-0000-0000-000000000004',
        '00000000-0000-0000-0000-000000000004',
        NULL,
        'Submitted',
        'Request created',
        NOW() - INTERVAL '6 hours'
    ),
    -- Log for Request 5: Created as Submitted
    (
        gen_random_uuid(),
        'a0000000-0000-0000-0000-000000000005',
        '00000000-0000-0000-0000-000000000003',
        NULL,
        'Submitted',
        'Request created',
        NOW() - INTERVAL '1 hour'
    ),
    -- Log for Request 6: Created, then Approved by admin
    (
        gen_random_uuid(),
        'a0000000-0000-0000-0000-000000000006',
        '00000000-0000-0000-0000-000000000004',
        NULL,
        'Submitted',
        'Request created',
        NOW() - INTERVAL '5 days'
    ),
    (
        gen_random_uuid(),
        'a0000000-0000-0000-0000-000000000006',
        '00000000-0000-0000-0000-000000000001',
        'Submitted',
        'Approved',
        'Disetujui. Tim electrical sudah dijadwalkan.',
        NOW() - INTERVAL '2 hours'
    )
ON CONFLICT DO NOTHING;

