
with weighed_random_shuffle_quiz_ids as (
	select quiz_id, 
		row_number() over (partition by quiz_id order by ts desc) as rn,
		RANDOM()*(case when res = false then 5.0 else 1.0 end) as chance
	from quiz_results
), next_quiz_id as (
	select quiz_id 
	from weighed_random_shuffle_quiz_ids
	where rn = 1
	order by chance desc
	limit 1
) --select * from next_quiz_id
 
select id, question, candidate_answers, correct_answer 
from (
	SELECT id, question, candidate_answers, correct_answer
	FROM quiz_data 
	WHERE $1 <@ tags) a --'{"S3"}'
	join next_quiz_id nqi on a.id = nqi.quiz_id
