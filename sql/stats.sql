with quiz_tags as (
	select id as qid, lower(unnest(tags)) as tag from quiz_data
), quiz_tags_service_map as (
	select * from quiz_tags join topic_service_map tsm 
    on tag = tsm.service
    where topic = 'basic'
), service_quiz_question_counts as (
	select service, count(distinct qid) as n from quiz_tags_service_map group by service
), min_date_minus_one as (
	select date(min(ts))-1 as d from quiz_results
), pre_init as (
	select d, qid, false as res from quiz_tags_service_map cross join min_date_minus_one
), user_res as (
    select r.quiz_id, r.ts, r.res as res
    from quiz_results r inner join quiz_data d
    on r.quiz_id = d.id
), user_res_enriched as (
    select date(ts) as d, qid, res 
    from user_res ur inner join quiz_tags_service_map qtsm
    on ur.quiz_id = qtsm.qid
    where topic = 'basic'
), user_res_plus_pre_init as (
    (select * from pre_init) union all (select * from user_res_enriched)
), ful as (
	select *, lag(res, 1) over (partition by qid order by d) as prev 
	from user_res_plus_pre_init order by d, qid
), filled as (
	select 
	d, qid, res, prev,
	CASE 
		WHEN res is NULL THEN coalesce(prev, FALSE)
	ELSE coalesce(res, FALSE) END as rfill
	from ful
	order by d, qid
), counts as (
    select d,
	SUM(CASE WHEN rfill = TRUE THEN 1 ELSE 0 END) AS n_correct,
    's3' as service
	from filled
	group by d
), counts_enriched as (
    select * from counts join service_quiz_question_counts using (service)
) select 
    d,
    service, 
    n_correct::float/n::float as pct
	from counts_enriched