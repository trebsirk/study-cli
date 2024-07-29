INSERT INTO
    quiz_data (question, candidate_answers, correct_answer, tags)
VALUES
    ($1, $2, $3, $4)
ON CONFLICT (question) DO NOTHING;