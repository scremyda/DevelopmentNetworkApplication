INSERT INTO public.autoparts(
    id, created_at, updated_at, deleted_at, name, description, brand, models, year, image, is_delete, user_id, status, price
)
VALUES (
           10,
           current_timestamp,
           current_timestamp,
           NULL,
           'Передний мотор Tesla Model E',
           'Передний мотор, Tesla Model У, 212098000S',
           'Tesla',
           'Tesla Model E',
           2020,
           'image5.jpg',
           false,
           3,
           'Available',
           50700
       );