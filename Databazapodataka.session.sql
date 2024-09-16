SELECT u.id, u.name, u.surname, a.date from users u right JOIN  attendance a on u.id = a.userid WHERE a.date BETWEEN '2024-07-01' AND '2024-09-01';

select CURRENT_DATE, CURRENT_DATE - 30;
