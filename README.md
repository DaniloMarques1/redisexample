The idea is that we have a Configuration entity in which you can define some default values such as timeout and a label.
You can only have one configuration definition, that way you can do one post (the first one) and the others post you do
will just update the existing row in the database. So whenever you do a get, the first one will take from the database
and save it in a cache so if you do a get one more time it will not get from the database, instead will get from the cache.
