

with service_quiz_ids as (
	select * from 
	(select id as qid, lower(unnest(tags)) as tag from quiz_data) a
	where tag = 's3'
), unique_quiz_ids as (
	select distinct(qid) from service_quiz_ids
), dday_minus_one as (
	select * from 
		(select date(date(min(ts))-1) as d, false as res from quiz_results) a
		cross join unique_quiz_ids
), user_res as (
	select * from quiz_results qr --quiz_id as qid, ts, res
	right outer join (
		select * from (
			select distinct(date(ts)) as dd from quiz_results
		) a
		cross join (select distinct(quiz_id) as qid from quiz_results) b
	) uqi 
	on qr.quiz_id = uqi.qid and date(qr.ts) = uqi.dd
	order by dd, qid
), user_res_last_daily_attempt as (
		select ts, qid, res, dd from 
			(select *, row_number() over (partition by dd, qid order by ts desc) as k
			from user_res) a
		where a.k = 1
), user_res_denullified as (
	select ts, case when res is NULL then FALSE else TRUE end as res, dd, qid 
	from user_res_last_daily_attempt
), filled_invalid as (
	(select *, lag(res, 1) over (partition by qid order by dd asc) as prev
			from user_res_denullified
			order by dd, qid, ts)
), filled_valid as (
	select  
			dd, ts, qid, res, prev,
		CASE 
			WHEN (res is NULL)
				THEN coalesce(prev, FALSE)
			WHEN prev is TRUE THEN TRUE
			ELSE coalesce(res, FALSE) 
			END as rfill
		from 
			filled_invalid
		order by dd, qid
), agg as (
	select dd, 'S3', 100*avg(case when rfill = TRUE then 1 else 0 end) as pct_complete from filled_valid
	group by 1 order by 1
) select * from agg