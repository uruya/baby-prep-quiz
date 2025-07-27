CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    category VARCHAR(50),
    question TEXT NOT NULL,
    options JSONB NOT NULL,
    correct_answer INT NOT NULL,
    explanation TEXT
);

INSERT INTO questions (id, category, question, options, correct_answer, explanation) VALUES
(1, 'pregnancy', '妊娠中の女性の体重増加の理想的な範囲は？',
    '["5kg未満","7kg〜12kg","15kg〜20kg","制限なし"]',
    1,
    '日本産科婦人科学会によると、妊娠中の理想的な体重増加は7kg〜12kgとされています。'),
(2, 'pregnancy', 'つわりが一般的に始まるのは妊娠何週目頃？',
    '["妊娠2週目頃","妊娠6週目頃","妊娠12週目頃","妊娠20週目頃"]',
    1,
    'つわりは一般的に妊娠6週目頃から始まり、16週目頃までに収まることが多いです。'),
(3, 'pregnancy', '妊娠中のパートナーへのサポートとして適切でないものは？',
    '["定期的な声かけと気遣い","家事の分担","栄養バランスの良い食事の準備","ストレス解消のためのお酒の提供"]',
    3,
    '妊娠中の飲酒は胎児に悪影響を与える可能性があるため、パートナーも一緒に控えることが望ましいです。');