-- select
--        r.id, r.room_name
-- from room r
-- where
--       r.id not in (select rr.room_id from room_restriction rr where "2022-08-08" < rr.check_out_date and "2022-09-09" > rr.check_in_date);
select * from reservation;