-- Development seed for the local tagma_test database.
-- Safe to run repeatedly: fixed seed records are updated instead of duplicated.

BEGIN;

INSERT INTO cities (id, name_ru, name_tm, code, status, created_at) VALUES
    (1, 'Ашхабад', 'Asgabat', 'ASB', true, now()),
    (2, 'Аркадаг', 'Arkadag', 'ARK', true, now())
ON CONFLICT (id) DO UPDATE SET
    name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm,
    code = EXCLUDED.code, status = EXCLUDED.status;

INSERT INTO districts (id, name_ru, name_tm, city_id, status, updated_at) VALUES
    (1, 'Беркарарлык', 'Berkararlyk', 1, true, now()),
    (2, 'Копетдаг', 'Kopetdag', 1, true, now()),
    (3, 'Аркадаг', 'Arkadag', 2, true, now())
ON CONFLICT (id) DO UPDATE SET
    name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm,
    city_id = EXCLUDED.city_id, status = EXCLUDED.status;

INSERT INTO days (id, name_ru, name_tm) VALUES
    (1, 'Понедельник', 'Duşenbe'),
    (2, 'Вторник', 'Sişenbe'),
    (3, 'Среда', 'Çarşenbe'),
    (4, 'Четверг', 'Penşenbe'),
    (5, 'Пятница', 'Anna'),
    (6, 'Суббота', 'Şenbe'),
    (7, 'Воскресенье', 'Ýekşenbe')
ON CONFLICT (id) DO UPDATE SET name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm;

INSERT INTO channel_types (id, name_ru, name_tm) VALUES
    (1, 'Магазин', 'Dükan'),
    (2, 'Супермаркет', 'Supermarket')
ON CONFLICT (id) DO UPDATE SET name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm;

INSERT INTO channel_structures (id, name_ru, name_tm) VALUES
    (1, 'Розничный', 'Bölek satuw'),
    (2, 'Оптовый', 'Lomaý satuw')
ON CONFLICT (id) DO UPDATE SET name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm;

INSERT INTO channel_sizes (id, name_ru, name_tm) VALUES
    (1, 'Малый', 'Kiçi'),
    (2, 'Средний', 'Orta')
ON CONFLICT (id) DO UPDATE SET name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm;

INSERT INTO channel_managements (id, name_ru, name_tm) VALUES
    (1, 'Частный', 'Hususy'),
    (2, 'Сетевой', 'Torlaýyn')
ON CONFLICT (id) DO UPDATE SET name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm;

INSERT INTO trade_channels (id, type_id, structure_id, size_id, management_id, status, updated_at) VALUES
    (1, 1, 1, 1, 1, true, now()),
    (2, 2, 1, 2, 2, true, now())
ON CONFLICT (id) DO UPDATE SET
    type_id = EXCLUDED.type_id, structure_id = EXCLUDED.structure_id,
    size_id = EXCLUDED.size_id, management_id = EXCLUDED.management_id,
    status = EXCLUDED.status;

INSERT INTO trade_categories (id, name_ru, name_tm, min_sum, max_sum, status, updated_at) VALUES
    (1, 'Стандарт', 'Standart', 0, 1000, true, now()),
    (2, 'Ключевой клиент', 'Esasy müşderi', 1000, 100000, true, now())
ON CONFLICT (id) DO UPDATE SET
    name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm,
    min_sum = EXCLUDED.min_sum, max_sum = EXCLUDED.max_sum,
    status = EXCLUDED.status;

INSERT INTO brands (id, name) VALUES
    (1, 'Tagma Water'),
    (2, 'Tagma Refreshments')
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;

INSERT INTO product_types (id, code, name_ru, name_tm, volume, count, status, updated_at) VALUES
    (1, 'WATER-05', 'Вода 0.5 л', 'Suw 0.5 l', 0.5, 12, true, now()),
    (2, 'JUICE-10', 'Сок 1 л', 'Şire 1 l', 1.0, 6, true, now()),
    (3, 'LEMONADE-05', 'Лимонад 0.5 л', 'Limonad 0.5 l', 0.5, 12, true, now())
ON CONFLICT (id) DO UPDATE SET
    code = EXCLUDED.code, name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm,
    volume = EXCLUDED.volume, count = EXCLUDED.count, status = EXCLUDED.status;

INSERT INTO volumes (id, volume) VALUES (1, 0.5), (2, 1.0)
ON CONFLICT (id) DO UPDATE SET volume = EXCLUDED.volume;

INSERT INTO products (id, code, taste_ru, taste_tm, price, d_price, status, updated_at, brand_id, product_type_id) VALUES
    (1, 'P-001', 'Без газа', 'Gazsyz', 3.50, 3.00, true, now(), 1, 1),
    (2, 'P-002', 'Газированная', 'Gazly', 4.00, 3.50, true, now(), 1, 1),
    (3, 'P-003', 'Лимон', 'Limon', 4.50, 4.00, true, now(), 1, 1),
    (4, 'P-004', 'Апельсин', 'Pyrtykal', 12.00, 10.00, true, now(), 2, 2),
    (5, 'P-005', 'Яблоко', 'Alma', 11.50, 9.50, true, now(), 2, 2),
    (6, 'P-006', 'Мультифрукт', 'Multimiwe', 13.00, 11.00, true, now(), 2, 2),
    (7, 'P-007', 'Кола', 'Kola', 6.00, 5.00, true, now(), 2, 3),
    (8, 'P-008', 'Лимон-лайм', 'Limon-laým', 6.00, 5.00, true, now(), 2, 3),
    (9, 'P-009', 'Ягодный', 'Miwe ir-iýmişli', 6.50, 5.50, true, now(), 2, 3)
ON CONFLICT (id) DO UPDATE SET
    code = EXCLUDED.code, taste_ru = EXCLUDED.taste_ru, taste_tm = EXCLUDED.taste_tm,
    price = EXCLUDED.price, d_price = EXCLUDED.d_price, status = EXCLUDED.status,
    brand_id = EXCLUDED.brand_id, product_type_id = EXCLUDED.product_type_id;

INSERT INTO drivers (id, name_ru, name_tm, number, status, updated_at) VALUES
    (1, 'Тестовый водитель', 'Synag sürüjisi', '+99360000001', true, now())
ON CONFLICT (id) DO UPDATE SET
    name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm,
    number = EXCLUDED.number, status = EXCLUDED.status;

INSERT INTO ekspeditors (id, name_ru, name_tm, number, status, updated_at) VALUES
    (1, 'Тестовый экспедитор', 'Synag ekspeditory', '+99360000002', true, now())
ON CONFLICT (id) DO UPDATE SET
    name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm,
    number = EXCLUDED.number, status = EXCLUDED.status;

INSERT INTO trade_agents (id, code, name_ru, name_tm, number, status, updated_at) VALUES
    (1, 'A-001', 'Тестовый агент', 'Synag agenti', '+99360000003', true, now())
ON CONFLICT (id) DO UPDATE SET
    code = EXCLUDED.code, name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm,
    number = EXCLUDED.number, status = EXCLUDED.status;

INSERT INTO agents_drivers (trade_agent_id, driver_id) VALUES (1, 1)
ON CONFLICT DO NOTHING;

INSERT INTO agents_ekspeditors (trade_agent_id, ekspeditor_id) VALUES (1, 1)
ON CONFLICT DO NOTHING;

INSERT INTO kinds (id, name_ru, name_tm) VALUES
    (1, 'Владелец', 'Eýesi'),
    (2, 'Менеджер', 'Dolandyryjy')
ON CONFLICT (id) DO UPDATE SET name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm;

INSERT INTO posts (id, name_ru, name_tm) VALUES
    (1, 'Владелец', 'Eýesi'),
    (2, 'Продавец', 'Satyjy')
ON CONFLICT (id) DO UPDATE SET name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm;

INSERT INTO trade_points (
    id, code, name_ru, name_tm, location_ru, location_tm, orientir_ru, orientir_tm,
    status, city_id, district_id, trade_agent_id, trade_channel_id,
    trade_category_id, day_id, ekspeditor_id, updated_at
) VALUES
    (1, 'TP-001', 'Магазин у парка', 'Seýilgähiň ýanyndaky dükan',
     'ул. Махтумкули, 10', 'Magtymguly köçesi, 10', 'Рядом с парком', 'Seýilgähiň ýanynda',
     true, 1, 1, 1, 1, 1, EXTRACT(ISODOW FROM CURRENT_DATE)::int, 1, now()),
    (2, 'TP-002', 'Супермаркет Центр', 'Merkez supermarketi',
     'пр. Битараплык, 25', 'Bitaraplyk şaýoly, 25', 'Первый этаж', 'Birinji gat',
     true, 1, 2, 1, 2, 2, EXTRACT(ISODOW FROM CURRENT_DATE)::int, 1, now()),
    (3, 'TP-003', 'Магазин Аркадаг', 'Arkadag dükany',
     'ул. Керим Гурбаннепесов, 5', 'Kerim Gurbannepesow köçesi, 5', 'У школы', 'Mekdebiň ýanynda',
     true, 2, 3, 1, 1, 1, (EXTRACT(ISODOW FROM CURRENT_DATE)::int % 7) + 1, 1, now()),
    (4, 'TP-004', 'Маркет Беркарар', 'Berkarar market',
     'ул. Ататюрка, 18', 'Atatürk köçesi, 18', 'Напротив торгового центра', 'Söwda merkeziniň garşysynda',
     true, 1, 1, 1, 2, 2, 2, 1, now()),
    (5, 'TP-005', 'Магазин Весна', 'Bahar dükany',
     'ул. Гёроглы, 42', 'Görogly köçesi, 42', 'Рядом с аптекой', 'Dermanhananyň ýanynda',
     true, 1, 1, 1, 1, 1, 3, 1, now()),
    (6, 'TP-006', 'Супермаркет Олимп', 'Olimp supermarketi',
     'пр. Туркменбаши, 77', 'Türkmenbaşy şaýoly, 77', 'Возле стадиона', 'Stadionyň ýanynda',
     true, 1, 2, 1, 2, 2, 4, 1, now()),
    (7, 'TP-007', 'Маркет Ак Алтын', 'Ak Altyn market',
     'ул. 10 лет Благополучия, 9', '10 ýyl Abadançylyk köçesi, 9', 'Первый поворот', 'Birinji öwrüm',
     true, 1, 2, 1, 1, 1, 5, 1, now()),
    (8, 'TP-008', 'Магазин Гюнеш', 'Güneş dükany',
     'ул. А.Ниязова, 31', 'A.Nyýazow köçesi, 31', 'За школой', 'Mekdebiň aňyrsynda',
     true, 1, 1, 1, 1, 1, 6, 1, now()),
    (9, 'TP-009', 'Супермаркет Шахер', 'Şäher supermarketi',
     'пр. Арчабил, 55', 'Arçabil şaýoly, 55', 'У кругового перекрёстка', 'Aýlawly çatrygyň ýanynda',
     true, 1, 2, 1, 2, 2, 7, 1, now()),
    (10, 'TP-010', 'Маркет Аркадаг Центр', 'Arkadag Merkez market',
     'ул. Акхан, 12', 'Akhan köçesi, 12', 'Рядом с автовокзалом', 'Awtomenziliň ýanynda',
     true, 2, 3, 1, 2, 2, 1, 1, now())
ON CONFLICT (id) DO UPDATE SET
    code = EXCLUDED.code, name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm,
    location_ru = EXCLUDED.location_ru, location_tm = EXCLUDED.location_tm,
    orientir_ru = EXCLUDED.orientir_ru, orientir_tm = EXCLUDED.orientir_tm,
    status = EXCLUDED.status, city_id = EXCLUDED.city_id, district_id = EXCLUDED.district_id,
    trade_agent_id = EXCLUDED.trade_agent_id, trade_channel_id = EXCLUDED.trade_channel_id,
    trade_category_id = EXCLUDED.trade_category_id, day_id = EXCLUDED.day_id,
    ekspeditor_id = EXCLUDED.ekspeditor_id;

INSERT INTO contacts (id, name_ru, name_tm, number, kind_id, post_id, trade_point_id) VALUES
    (1, 'Анна', 'Anna', '+99361000001', 1, 1, 1),
    (2, 'Мердан', 'Merdan', '+99361000002', 2, 2, 2),
    (3, 'Бегенч', 'Begenç', '+99361000003', 1, 1, 3),
    (4, 'Майса', 'Maýsa', '+99361000004', 2, 2, 4),
    (5, 'Сердар', 'Serdar', '+99361000005', 1, 1, 5),
    (6, 'Айна', 'Aýna', '+99361000006', 2, 2, 6),
    (7, 'Ровшен', 'Röwşen', '+99361000007', 1, 1, 7),
    (8, 'Лейли', 'Leýli', '+99361000008', 2, 2, 8),
    (9, 'Довлет', 'Döwlet', '+99361000009', 1, 1, 9),
    (10, 'Марал', 'Maral', '+99361000010', 2, 2, 10)
ON CONFLICT (id) DO UPDATE SET
    name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm,
    number = EXCLUDED.number, kind_id = EXCLUDED.kind_id,
    post_id = EXCLUDED.post_id, trade_point_id = EXCLUDED.trade_point_id;

INSERT INTO routes (id, trade_point_id, trade_agent_id, ekspeditor_id, day_id, status, updated_at) VALUES
    (1, 1, 1, 1, EXTRACT(ISODOW FROM CURRENT_DATE)::int, true, now()),
    (2, 2, 1, 1, EXTRACT(ISODOW FROM CURRENT_DATE)::int, true, now()),
    (3, 3, 1, 1, (EXTRACT(ISODOW FROM CURRENT_DATE)::int % 7) + 1, true, now()),
    (4, 4, 1, 1, 2, true, now()),
    (5, 5, 1, 1, 3, true, now()),
    (6, 6, 1, 1, 4, true, now()),
    (7, 7, 1, 1, 5, true, now()),
    (8, 8, 1, 1, 6, true, now()),
    (9, 9, 1, 1, 7, true, now()),
    (10, 10, 1, 1, 1, true, now())
ON CONFLICT (id) DO UPDATE SET
    trade_point_id = EXCLUDED.trade_point_id, trade_agent_id = EXCLUDED.trade_agent_id,
    ekspeditor_id = EXCLUDED.ekspeditor_id, day_id = EXCLUDED.day_id, status = EXCLUDED.status;

DELETE FROM route_orders WHERE trade_agent_id = 1;
INSERT INTO route_orders (order_number, trade_point_id, trade_agent_id, day_id) VALUES
    (1, 1, 1, EXTRACT(ISODOW FROM CURRENT_DATE)::int),
    (2, 2, 1, EXTRACT(ISODOW FROM CURRENT_DATE)::int),
    (1, 3, 1, (EXTRACT(ISODOW FROM CURRENT_DATE)::int % 7) + 1),
    (1, 4, 1, 2),
    (1, 5, 1, 3),
    (1, 6, 1, 4),
    (1, 7, 1, 5),
    (1, 8, 1, 6),
    (1, 9, 1, 7),
    (1, 10, 1, 1);

INSERT INTO payment_types (id, name_ru, name_tm) VALUES
    (1, 'Наличные', 'Nagt'),
    (2, 'Терминал', 'Terminal'),
    (3, 'Перевод', 'Geçirim'),
    (4, 'Промо', 'Promo')
ON CONFLICT (id) DO UPDATE SET name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm;

-- ID 1 is the application's conventional "no discount" record.
INSERT INTO discounts (
    id, name_ru, name_tm, start, "end", sort_order, status,
    is_procent, is_direct, is_increase, is_group, created_at
) VALUES
    (1, 'Без скидки', 'Arzanladyşsyz', CURRENT_DATE - INTERVAL '1 year',
     CURRENT_DATE + INTERVAL '10 years', 0, true, false, true, false, false, now())
ON CONFLICT (id) DO UPDATE SET
    name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm,
    start = EXCLUDED.start, "end" = EXCLUDED."end", status = EXCLUDED.status;

INSERT INTO orders (
    id, trade_point_id, trade_agent_id, sum, sale_sum, created_at,
    status, is_closed, is_credit, payment_type_id
) VALUES
    (1, 1, 1, 47.00, 47.00, CURRENT_DATE, true, false, false, 1),
    (2, 2, 1, 60.00, 60.00, CURRENT_DATE, true, false, true, 2),
    (3, 1, 1, 24.00, 24.00, CURRENT_DATE - INTERVAL '1 day', true, true, false, 1)
ON CONFLICT (id) DO UPDATE SET
    trade_point_id = EXCLUDED.trade_point_id, trade_agent_id = EXCLUDED.trade_agent_id,
    sum = EXCLUDED.sum, sale_sum = EXCLUDED.sale_sum, created_at = EXCLUDED.created_at,
    status = EXCLUDED.status, is_closed = EXCLUDED.is_closed,
    is_credit = EXCLUDED.is_credit, payment_type_id = EXCLUDED.payment_type_id;

INSERT INTO order_lists (
    id, product_id, discount_id, order_id, price, sale_price, coefficient, count, status
) VALUES
    (1, 1, 1, 1, 35.00, 35.00, 0, 10, true),
    (2, 4, 1, 1, 12.00, 12.00, 0, 1, true),
    (3, 7, 1, 2, 60.00, 60.00, 0, 10, true),
    (4, 4, 1, 3, 24.00, 24.00, 0, 2, true)
ON CONFLICT (id) DO UPDATE SET
    product_id = EXCLUDED.product_id, discount_id = EXCLUDED.discount_id,
    order_id = EXCLUDED.order_id, price = EXCLUDED.price, sale_price = EXCLUDED.sale_price,
    coefficient = EXCLUDED.coefficient, count = EXCLUDED.count, status = EXCLUDED.status;

INSERT INTO payments (
    id, order_id, currency, payment_type_id, created_at,
    trade_agent_id, status, count, is_agent
) VALUES
    (1, 1, 20.00, 1, now(), 1, false, 0, false),
    (2, 2, 25.00, 2, now(), 1, false, 0, false),
    (3, 3, 24.00, 1, CURRENT_DATE - INTERVAL '1 day', 1, true, 0, false)
ON CONFLICT (id) DO UPDATE SET
    order_id = EXCLUDED.order_id, currency = EXCLUDED.currency,
    payment_type_id = EXCLUDED.payment_type_id, trade_agent_id = EXCLUDED.trade_agent_id,
    status = EXCLUDED.status, count = EXCLUDED.count, is_agent = EXCLUDED.is_agent;

INSERT INTO debts (id, trade_agent_id, total) VALUES (1, 1, 5.00)
ON CONFLICT (id) DO UPDATE SET trade_agent_id = EXCLUDED.trade_agent_id, total = EXCLUDED.total;

INSERT INTO dealers (
    id, name_ru, name_tm, address_ru, address_tm, company, phone_number,
    email, city_id, account, debt, created_at, updated_at, status
) VALUES
    (1, 'Тестовый дилер', 'Synag dileri', 'Ашхабад', 'Asgabat', 'Tagma Test',
     '+99362000001', 'dealer@example.test', 1, 1000.00, 100.00, now(), now(), true)
ON CONFLICT (id) DO UPDATE SET
    name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm,
    company = EXCLUDED.company, city_id = EXCLUDED.city_id,
    account = EXCLUDED.account, debt = EXCLUDED.debt, status = EXCLUDED.status;

INSERT INTO gift_types (id, name_ru, name_tm) VALUES
    (1, 'Товар', 'Haryt'),
    (2, 'Сувенир', 'Ýadygärlik')
ON CONFLICT (id) DO UPDATE SET name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm;

INSERT INTO furnitures (id, name_ru, name_tm, created_at, updated_at) VALUES
    (1, 'Холодильник', 'Sowadyjy', now(), now()),
    (2, 'Стеллаж', 'Tekje', now(), now())
ON CONFLICT (id) DO UPDATE SET name_ru = EXCLUDED.name_ru, name_tm = EXCLUDED.name_tm;

INSERT INTO trade_points_furnitures (id, trade_point_id, furniture_id, count, created_at) VALUES
    (1, 1, 1, 1, now()),
    (2, 1, 2, 2, now()),
    (3, 2, 2, 4, now()),
    (4, 4, 1, 2, now()),
    (5, 5, 2, 2, now()),
    (6, 6, 1, 3, now()),
    (7, 7, 2, 3, now()),
    (8, 9, 1, 2, now()),
    (9, 10, 2, 4, now())
ON CONFLICT (id) DO UPDATE SET
    trade_point_id = EXCLUDED.trade_point_id,
    furniture_id = EXCLUDED.furniture_id, count = EXCLUDED.count;

-- Legacy SHA-1 hashes are intentionally used here. On first successful login,
-- the application upgrades each password to bcrypt.
INSERT INTO users (id, name_tm, name_ru, username, password, type, status, created_at) VALUES
    (1, 'Administrator', 'Администратор', 'admin',
     '61736c6b64676f69616c736b6466616c7364666c616b73646a6866f865b53623b121fd34ee5426c792e5c33af8c227',
     'admin', true, now()),
    (2, 'Synag agenti', 'Тестовый агент', 'agent',
     '61736c6b64676f69616c736b6466616c7364666c616b73646a6866abcd42066f55af2adfdfc8c4ba0ad5636fe0af29',
     'agent', true, now()),
    (3, 'Synag ekspeditory', 'Тестовый экспедитор', 'ekspeditor',
     '61736c6b64676f69616c736b6466616c7364666c616b73646a686617d260a7b0b00945876467d68c21860191ed0896',
     'ekspeditor', true, now())
ON CONFLICT (id) DO UPDATE SET
    name_tm = EXCLUDED.name_tm, name_ru = EXCLUDED.name_ru,
    username = EXCLUDED.username, password = EXCLUDED.password,
    type = EXCLUDED.type, status = EXCLUDED.status;

INSERT INTO user_trade_agents (id, user_id, trade_agent_id) VALUES (1, 2, 1)
ON CONFLICT (id) DO UPDATE SET user_id = EXCLUDED.user_id, trade_agent_id = EXCLUDED.trade_agent_id;

INSERT INTO user_ekspeditors (id, user_id, ekspeditor_id) VALUES (1, 3, 1)
ON CONFLICT (id) DO UPDATE SET user_id = EXCLUDED.user_id, ekspeditor_id = EXCLUDED.ekspeditor_id;

-- Keep every owned sequence ahead of explicitly assigned seed values.
DO $$
DECLARE
    sequence_record record;
    max_value bigint;
BEGIN
    FOR sequence_record IN
        SELECT
            table_namespace.nspname AS table_schema,
            table_class.relname AS table_name,
            table_attribute.attname AS column_name,
            sequence_class.oid::regclass AS sequence_name
        FROM pg_class sequence_class
        JOIN pg_depend dependency
          ON dependency.objid = sequence_class.oid
         AND dependency.deptype IN ('a', 'i')
        JOIN pg_class table_class
          ON table_class.oid = dependency.refobjid
        JOIN pg_namespace table_namespace
          ON table_namespace.oid = table_class.relnamespace
        JOIN pg_attribute table_attribute
          ON table_attribute.attrelid = table_class.oid
         AND table_attribute.attnum = dependency.refobjsubid
        WHERE sequence_class.relkind = 'S'
          AND table_namespace.nspname = 'public'
    LOOP
        EXECUTE format(
            'SELECT MAX(%I) FROM %I.%I',
            sequence_record.column_name,
            sequence_record.table_schema,
            sequence_record.table_name
        ) INTO max_value;

        IF max_value IS NOT NULL THEN
            PERFORM setval(sequence_record.sequence_name, max_value, true);
        END IF;
    END LOOP;
END
$$;

COMMIT;
