
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
		select ts, dd, qid, res from 
			(select *, row_number() over (partition by dd, qid order by ts desc) as k
			from user_res) a
		where a.k = 1
), x as (
	select *,
		first_value(res) OVER (PARTITION BY res_partition order by ts ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) as last_val
	from
	(select *, count(res) over (partition by qid order by dd) as res_partition
	from user_res_last_daily_attempt) a
	order by qid, dd
), y as (
	select *, case when res_partition = 0 then false else coalesce(res, last_val) end as last_val2 from x
), agg as (
	select dd, 100*avg(case when last_val2 is TRUE then 1 else 0 end) as pct_complete 
	from y
	group by 1 order by 1
)
-- select qid, last_value(res) over (
-- 	partition by qid order by dd ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
-- 	) as last_val 
-- from user_res_last_daily_attempt order by qid --group by 1 order by qid
select dd, 'S3', pct_complete from agg
-- select * from user_res where qid = 4
-- select * from user_res_last_daily_attempt where qid = 10
-- select dd, qid, res, prev from filled_invalid where qid = 10 --dd >= '2024-07-01' and
-- select dd, qid, res, prev, cum_res, cum_res2 from filled_valid2 where qid = 11
-- select dd, qid, res, last_val, cum from cum_res_cte where qid = 11
-- select dd, qid, res, res_partition, last_val, last_val2 from y order by qid,dd-- where qid = 11 or qid = 3

--select * from filled_valid where qid = 10
--select * from quiz_data where id = 11
--select dd, count(distinct(qid)) from user_res group by 1 order by 1qid = 4
-- select * from agg
--select * from filled_valid where qid = 10
--select * from quiz_data where id = 11
--select dd, count(distinct(qid)) from user_res group by 1 order by 1