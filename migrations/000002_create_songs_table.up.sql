CREATE TABLE songs (
    id SERIAL PRIMARY KEY,          
    title VARCHAR(255) NOT NULL,
    group_name VARCHAR(255) NOT NULL,
    release_date DATE,              
    lyrics TEXT,                   
    link VARCHAR(255)
);

CREATE INDEX idx_songs_title ON songs(title);
CREATE INDEX idx_songs_group_name ON songs(group_name);
CREATE INDEX idx_songs_release_date ON songs(release_date);
CREATE INDEX idx_songs_lyrics ON songs USING GIN (to_tsvector('english', lyrics));
CREATE INDEX idx_songs_link ON songs(link);

INSERT INTO songs (title, group_name, release_date, lyrics, link) VALUES
-- Kaleo
('Way Down We Go', 'Kaleo', '2016-01-29', 
'Oh, Father tell me, do we get what we deserve?
Whoa, we get what we deserve

And way down we go
Way down we go
Say way down we go
Way down we go

You let your feet run wild
Time has come as we all, oh, go down
Yeah but for the fall, ooh, my
Do you dare to look him right in the eyes? Yeah

Oh, ''cause they will run you down, down ''til the dark
Yes and they will run you down, down ''til you fall
And they will run you down, down ''til you go
Yeah, so you can''t crawl no more

And way down we go
Way down we go
Say way down we go

Oh, ''cause they will run you down, down ''til you fall
Way down we go

Oh baby, yeah
Oh, baby
Baby
Way down we go
Yeah

And way down we go
Way down we go
Say way down we go, ooh
Way down we go', 
'https://www.youtube.com/watch?v=0-7IHOXkiV8'),

('All the Pretty Girls', 'Kaleo', '2016-06-30', 
'All the pretty girls like Samuel
Oh he really doesn''t share
Though it''s more than he can handle
Life is anything but fair, life is anything but fair

Just as soon as they turn older
He''ll come and sweep them off their feet
It''s only making me feel smaller
All the hidden love beneath

So won''t you lay me, won''t you lay me down
Won''t you lay me, won''t you lay me down
Won''t you lay me, won''t you lay me down
Won''t you lay me, won''t you lay....me down

All alone, alone again
No one lends a helping hand
I have waited, I have waited
Takes it''s toll, one''s foolish pride
How long before I see the light
I have waited, I have waited for you to lay me down

Sail on by, sail on by for now
They play naked in the water
You know it''s hard, heaven knows I''ve tried
But it just keeps getting harder

So won''t you lay me, won''t you lay me down
Won''t you lay me, won''t you lay me down
Won''t you lay me, won''t you lay me down
Won''t you lay me, won''t you lay....

Oh won''t you lay me, won''t you lay me down
Won''t you lay me, oh won''t you lay me down
Oh Won''t you lay me, say won''t you lay me down
Won''t you lay me down

Oh I''ll wait, I''ll wait, I''ll wait, I''ll wait for you
Yeah I''ll wait, I''ll wait, I''ll wait, I''ll wait for you
Oh I''ll wait, I''ll wait, I''ll wait, I''ll wait for you
Oh I''ll wait, I''ll wait, I''ll wait, I''ll wait for you

For you to lay me
Won''t you lay me down', 
'https://www.youtube.com/watch?v=FNwgOkl5nRY'),

('I can''t go on without you', 'Kaleo', '2016-03-12', 
'Well, they thought they were made for each other
Only thinking of one another
Never thinking just for one second
She would take a different attraction

We don''t want that
We don''t want that
We don''t want that, oh no
We don''t want that
We don''t want that
We don''t want that, oh no

I can''t go on without you
I can''t go on without you
Can''t go on without you, yeah
I can''t go on without you

Oh, so what''s the point of breaking my sweet heart?
She wanted me to let down my guard
Well, you know what they say, it''s
It''s better that way, so
So, you better hush and walk away

We don''t want that
We don''t want that
We don''t want that, oh no
We don''t want that
We don''t want that
We don''t want that, oh no

I can''t go on without you
I can''t go on without you, oh, Lord
Can''t go on without you
I can''t go on, won''t go on
Living on, without you

Oh, yeah
Woah, oh
Well, was I supposed to wait for you sweetheart?
And hide away the shame, yes I keep it all inside
Though the thought had crossed my mind
To do all the things I''ll regret, we don''t want that

We don''t want that
We don''t want that
We don''t want that, oh no
We don''t want that
We don''t want that
We don''t want that, oh no

I can''t, I can''t, I can''t go on without you
I can''t go on without you, oh, Lord
Go on without you
I can''t go on without you, babe

Yeah
Oh, she loves me
She loves me not
She loves me
My love don''t love me

Oh, so what is left but a broken man?
''Cause nothing hurts like a woman can

I can''t go on without you
I can''t go on without you, oh yeah
Can''t go on without you
I can''t go on without you, oh
Oh, without you Lord, without you
Without you, babe
Without you, oh
Oh, Lord
You', 
'https://www.youtube.com/watch?v=jfNOdsvMke4'),

-- Of Monsters and Men
('Little Talks', 'Of Monsters and Men', '2011-12-27', 
'Hey! Hey! Hey!
I don''t like walking around this old and empty house
So hold my hand, I''ll walk with you, my dear
The stairs creak as you sleep, it''s keeping me awake
It''s the house telling you to close your eyes

And some days I can''t even dress myself
It''s killing me to see you this way

''Cause though the truth may vary
This ship will carry our bodies safe to shore

Hey! Hey! Hey!

There''s an old voice in my head that''s holding me back
Well tell her that I miss our little talks
Soon it will be over and buried with our past
We used to play outside when we were young
And full of life and full of love.

Some days I don''t know if I am wrong or right
Your mind is playing tricks on you, my dear

''Cause though the truth may vary
This ship will carry our bodies safe to shore

Hey!
Don''t listen to a word I say
Hey!
The screams all sound the same
Hey!

Though the truth may vary
This ship will carry our bodies safe to shore

Hey!
Hey!

You''re gone, gone, gone away
I watched you disappear
All that''s left is a ghost of you.
Now we''re torn, torn, torn apart,
There''s nothing we can do
Just let me go we''ll meet again soon
Now wait, wait, wait for me
Please hang around
I''ll see you when I fall asleep

Hey!
Don''t listen to a word I say
Hey!
The screams all sound the same
Hey!
Though the truth may vary
This ship will carry our bodies safe to shore

Don''t listen to a word I say
Hey!
The screams all sound the same
Hey!

Though the truth may vary
This ship will carry our bodies safe to shore

Though the truth may vary
This ship will carry our bodies safe to shore

Though the truth may vary
This ship will carry our bodies safe to shore', 
'https://www.youtube.com/watch?v=IY8rOSyR5Rw'),

('Crystals', 'Of Monsters and Men', '2015-09-25', 
'Lost in skies of powdered gold
Caught in clouds of silver ropes
Showered by the empty hopes
As I tumble down, falling fast to the ground

I know I''ll wither so peel away the bark
''Cause nothing grows when it is dark
In spite of all my fears, I can see it all so clear
I see it all so clear

Whoa-o-o-o, cover your crystal eyes
And feel the tones that tremble down your spine
Whoa-o-o-o, cover your crystal eyes
And let your colors bleed and blend with mine

Making waves in pitch black sand
Feel the salt dance on my hands
Raw and charcoal colored thighs feel so cold
And my skin feels so paper-thin

I know I''ll wither so peel away the bark
''Cause nothing grows when it is dark
In spite of all my fears, I can see it all so clear
I see it all so clear

Whoa-o-o-o, cover your crystal eyes
And feel the tones that tremble down your spine
Whoa-o-o-o, cover your crystal eyes
And let your colors bleed and blend with mine

But I''m okay in see-through skin
I forgive what is within
''Cause I''m in this house
I''m in this home
All my time

Whoa-o-o-o, cover your crystal eyes
And feel the tones that tremble down your spine
Whoa-o-o-o, cover your crystal eyes
And let your colors bleed and blend with mine', 
'https://www.youtube.com/watch?v=_-PgPZ3F9P4'),

('Dirty Paws', 'Of Monsters and Men', '2011-12-27', 
'Jumping up and down the floor,
My head is an animal.
And once there was an animal,
It had a son that mowed the lawn.
The son was an OK guy,
They had a pet dragonfly.
The dragonfly, it ran away,
But it came back with a story to say.

Her dirty paws and furry coat,
She ran down the forest slopes.
The forest of talking trees,
They used to sing about the birds and the bees.
The bees had declared a war,
The sky wasn''t big enough for them all.
The birds, they got help from below,
From dirty paws and the creatures of snow.

So for a while things were cold,
They were scared down in their holes.
The forest that once was green
Was colored black by those killing machines.
But she and her furry friends
Took down the queen bee and her men.
And that''s how the story goes,
The story of the beast with those four dirty paws.', 
'https://www.youtube.com/watch?v=mCHUw7ACS8o'),

-- David Kushner
('Daylight', 'David Kushner', '2023-09-01', 
'Telling myself I won''t go there
Oh, but I know that I won''t care
Tryna wash away all the blood I''ve spilt
This lust is a burden that we both share
Two sinners can''t atone from a lone prayer
Souls tied, intertwined by pride and guilt

(Ooh) There''s darkness in the distance
From the way that I''ve been livin''
(Ooh) But I know I can''t resist it

Oh, I love it and I hate it at the same time
You and I drink the poison from the same vine
Oh, I love it and I hate it at the same time
Hidin'' all of our sins from the daylight
From the daylight, runnin'' from the daylight
From the daylight, runnin'' from the daylight
Oh, I love it and I hate it at the same time

Tellin'' myself it''s the last time
Can you spare any mercy that you might find
If I''m down on my knees again?
Deep down, way down, Lord, I try
Try to follow your light, but it''s nighttime
Please don''t leave me in the end

(Ooh) There''s darkness in the distance
I''m beggin'' for forgiveness
(Ooh) But I know I might resist it, oh

Oh, I love it and I hate it at the same time
You and I drink the poison from the same vine
Oh, I love it and I hate it at the same time
Hidin'' all of our sins from the daylight
From the daylight, runnin'' from the daylight
From the daylight, runnin'' from the daylight
Oh, I love it and I hate it at the same time
Oh, I love it and I hate it at the same time
You and I drink the poison from the same vine
Oh, I love it and I hate it at the same time
Hidin'' all of our sins from the daylight
From the daylight, runnin'' from the daylight
From the daylight, runnin'' from the daylight
Oh, I love it and I hate it at the same time', 
'https://www.youtube.com/watch?v=MoN9ql6Yymw'),

('Darkerside', 'David Kushner', '2024-08-31', 
'Been running too long
Trying to catch my breath
There''s a war up against
My heart and head
And there ain''t always
Blood in a fight
But you bring me back to the light
But you bring me back to the light

Oh why, oh why
Am I standing on the edge of my Darkerside
Oh I, oh I
I know it''s so wrong but it feels so right
Oh my, oh my
There''s so many things I''m tempted by
But you bring me back to the light
Oh my

I''m spiraling down and out of control
There''s a war no one sees inside my soul
The Lord exposes the deepest lies
And you bring me back to the light
Yeah you bring me back to the light

Oh why, oh why
Am I standing on the edge of my Darkerside
Oh I, oh I
I know it''s so wrong but it feels so right
Oh my, oh my
There''s so many things I''m tempted by
But you bring me back to the light

Yeah you bring me back to the light

But you bring me back to the light
Oh my', 
'https://www.youtube.com/watch?v=sDNM-kPa2jw'),

('Humankind', 'David Kushner', '2024-06-21', 
'I met the devil Sunday mornin'' with his hands in the air
He blew his paycheck on the plate that they were passin''
He''s testifyin'' in the light, but he was heartless in prayer
He had the spirit and expensive taste in fashion

I''m the one that you came and slaughtered
You spin me around
I was lookin'' for livin'' water
You just let me drown

I put my faith in a sinner''s town
Land of the free chained to the ground
When I look for kindness now
Humankind just lets me down

Oh, oh, oh-oh, oh
Oh, my heaven''s cussing me out
Oh, oh, oh-oh, oh
Humankind just lets me down

They lost the message in a bottle, now it''s covered in blood
They sell a savin'' just to make a Great Commission
We love the ones who hurt us, and we hurt the ones that we love
We''re sacrificin'' one another for tradition

I''m the one that you came and slaughtered
You spin me around
I was lookin'' for livin'' water
You just let me drown

I put my faith in a sinner''s town
Land of the free chained to the ground
When I look for kindness now
Humankind just lets me down

Oh, oh, oh-oh, oh
Oh, my heaven''s cussing me out
Oh, oh, oh-oh, oh
Humankind just lets me down

When it''s all said and done
I''m just a man, you''re just a woman
So take my hand
I''m only human
I need you to take me

Home
Home
Home
Humankind just lets me down', 
'https://www.youtube.com/watch?v=WXwKXouZHFw')