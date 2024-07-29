SELECT id, question, candidate_answers, correct_answer
FROM quiz_data 
WHERE tags && $1 
ORDER BY RANDOM() asc
LIMIT 1