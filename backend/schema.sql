CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    category VARCHAR(50),
    question TEXT NOT NULL,
    options JSONB NOT NULL,
    correct_answer INT NOT NULL,
    explanation TEXT
);

-- データが存在しない場合のみ挿入する
INSERT INTO questions (category, question, options, correct_answer, explanation) VALUES
('pregnancy', '妊娠中の女性の体重増加の理想的な範囲は？',
    '["5kg未満","7kg〜12kg","15kg〜20kg","制限なし"]',
    1,
    '日本産科婦人科学会によると、妊娠中の理想的な体重増加は7kg〜12kgとされています。')
ON CONFLICT (id) DO NOTHING;

INSERT INTO questions (category, question, options, correct_answer, explanation) VALUES
('pregnancy', 'つわりが一般的に始まるのは妊娠何週目頃？',
    '["妊娠2週目頃","妊娠6週目頃","妊娠12週目頃","妊娠20週目頃"]',
    1,
    'つわりは一般的に妊娠6週目頃から始まり、16週目頃までに収まることが多いです。')
ON CONFLICT (id) DO NOTHING;

INSERT INTO questions (category, question, options, correct_answer, explanation) VALUES
('pregnancy', '妊娠中のパートナーへのサポートとして適切でないものは？',
    '["定期的な声かけと気遣い","家事の分担","栄養バランスの良い食事の準備","ストレス解消のためのお酒の提供"]',
    3,
    '妊娠中の飲酒は胎児に悪影響を与える可能性があるため、パートナーも一緒に控えることが望ましいです。')
ON CONFLICT (id) DO NOTHING;

-- カテゴリ: birth (出産の準備)
INSERT INTO questions (id, category, question, options, correct_answer, explanation) VALUES
(4, 'birth', '陣痛が始まったと感じた時、まず何をすべき？',
    '["すぐに病院へ向かう","お風呂に入る","病院に電話して指示を仰ぐ","食事をとる"]',
    2,
    'まずは慌てずに病院へ電話し、現在の状況（陣痛の間隔など）を伝えて指示を仰ぐのが一般的です。')
ON CONFLICT (id) DO NOTHING;

INSERT INTO questions (id, category, question, options, correct_answer, explanation) VALUES
(5, 'birth', '入院バッグに出産後すぐ必要になるものは？',
    '["ベビー服","哺乳瓶","産褥ショーツ","離乳食"]',
    2,
    '産褥ショーツは、産後の悪露（おろ）に対応するために、出産後すぐに必要となるママのためのアイテムです。')
ON CONFLICT (id) DO NOTHING;

-- カテゴリ: baby-care (赤ちゃんのお世話)
INSERT INTO questions (id, category, question, options, correct_answer, explanation) VALUES
(6, 'baby-care', '新生児のおむつ替えの適切な頻度は？',
    '["1日に3回","授乳のたび、または汚れたらその都度","朝と夜の2回","泣いたときだけ"]',
    1,
    '新生児は排泄の回数が多いため、授乳のたびやおむつが汚れていることに気づいた都度、交換してあげるのが理想的です。')
ON CONFLICT (id) DO NOTHING;