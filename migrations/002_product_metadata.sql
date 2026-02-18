-- Таблица параметров изделий AMANI (универсальная для всех категорий)
CREATE TABLE IF NOT EXISTS product_metadata (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) NOT NULL,          -- Артикул (связь с МойСклад)
    category VARCHAR(30),              -- Платье, Брюки, Пальто и т.д.
    size_name VARCHAR(10) NOT NULL,    -- Размер (S, M, 44, 46...)
    
    -- Основные обхваты изделия (см)
    bust_full NUMERIC(5,1),            -- Обхват груди
    waist_full NUMERIC(5,1),           -- Обхват талии
    hips_full NUMERIC(5,1),            -- Обхват бедер
    
    -- Дополнительные конструктивные параметры (см)
    shoulder_width NUMERIC(5,1),       -- Ширина плеч (важно для жакетов/пальто)
    sleeve_length NUMERIC(5,1),        -- Длина рукава
    product_length NUMERIC(5,1),       -- Длина изделия
    crotch_depth NUMERIC(5,1),         -- Высота сидения (важно для брюк)
    
    -- Характеристики ткани и посадки
    silhouette VARCHAR(30),            -- Приталенный, Прямой, Трапеция, Оверсайз
    stretch_level VARCHAR(20),         -- Нет, Слабая, Средняя, Высокая (для трикотажа)
    lining_type VARCHAR(50),           -- Тип подклада/утеплителя (для курток/пальто)
    
    -- Рекомендации системы
    rec_height_min INT,                -- Мин. рост (напр. 160)
    rec_height_max INT,                -- Макс. рост (напр. 175)
    ease_allowance_cm NUMERIC(4,1),    -- Заложенная свобода облегания (см)

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(sku, size_name)             -- Защита от дублирования размеров в одном артикуле
);

-- Таблица связки контента (Instagram/Ads) с товаром
CREATE TABLE IF NOT EXISTS content_map (
    shortcode VARCHAR(50) PRIMARY KEY, -- ID поста или спец. код ссылки
    product_sku VARCHAR(50) NOT NULL,  -- Артикул изделия
    campaign_name VARCHAR(100),        -- Название рекламы или имя блогера
    created_at TIMESTAMP DEFAULT NOW()
);

-- Индекс для быстрого поиска по артикулу
CREATE INDEX IF NOT EXISTS idx_product_metadata_sku ON product_metadata(sku);
