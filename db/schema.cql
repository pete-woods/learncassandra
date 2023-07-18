create keyspace if not exists example with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

create table if not exists example.tweet(timeline text, id UUID, text text, PRIMARY KEY(id));

create index if not exists on example.tweet(timeline);