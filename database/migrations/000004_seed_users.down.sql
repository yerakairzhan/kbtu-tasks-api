DELETE FROM user_friends
WHERE user_id IN (
  'user-01','user-02','user-03','user-04','user-05','user-06','user-07','user-08','user-09','user-10',
  'user-11','user-12','user-13','user-14','user-15','user-16','user-17','user-18','user-19','user-20'
);

DELETE FROM users
WHERE id IN (
  'user-01','user-02','user-03','user-04','user-05','user-06','user-07','user-08','user-09','user-10',
  'user-11','user-12','user-13','user-14','user-15','user-16','user-17','user-18','user-19','user-20'
);
