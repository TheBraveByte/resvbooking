select r.id,
       r.first_name,
       r.last_name,
       r.email,
       r.phone_number,
       r.room_id,
       r.check_in_date,
       r.check_in_date,
       r.updated_at,
       r.created_at,
       r.processed
from reservation r
         left join rooms rm on (r.room_id = rm.id)
where r.processed = 0
order by r.check_in_date