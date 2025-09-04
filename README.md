# Wengine the search engine: Blueprint
## Explore and Index Services
The Explore and Index Services could in principle work without the "User Web Interface",
as long as the PostgresDB holding the Indexed websites and the RedisDB "to index" queue, as well as
PageRank cache are online.

### Crawler
1. Written in Golang. Is in a Docker cluster managed by The Scaling Service.
2. Crawls pages: Gets the html, puts the page html, as well as all of the pages links, and image links and their number into the "to index" queue
3. Breadth first search.
4. Also puts a simplified version with "page path - list of connected pages - page rank" into
a Redis cache for the PageRank alg, and to know what has alredy been crawled. Uses PostgresDB directly as fallback.
5. Available in a Docker image form. To be run in a Docker cluster.

### Indexer
1. Written in Python. Is in a Docker cluster managed by The Scaling Service.
2. Takes from Crawler given "to index" Redis queue, counts words,
and extracts Links, Title and Content Snippet, puts into PostgresDB.
4. Available in a Docker image form. To be run in a Docker cluster.

### Page Rank
1. Written in Golang. Is in a Docker cluster managed by The Scaling Service.
2. Goes over the Redis entries prepared for it by the Indexer and updates their page rank.
3. May have to use the PostgresDB directly as fallback if its Redis cache gets too fat, or otherwise fails.
4. Once in a while updates the page rank of webpages in the PostgresDB
5. Available in a Docker image form. To be run in a Docker cluster.

---

## Databases
Both Postgres and Redis are used as an in-between.
Postgres is in-between Explore/Index Services and the User Web Interface.
Redis is used in-between the crawler and the indexer,
and as a PageRank cache, and to track alredy crawled pages.

### PostgresDB
1. Available in a Docker image form. To make it simple to run.
2. Hopefully will run it in the cloud for Demonstration purposes.

#### Database schema
The Created at and Updated at fields are used for debugging,
and to make it possible to re crawl, and re index pages after a certain period of time.
It may make more sense to add and use an additional, Refreshed at field to do that however.

##### **Pages Table**
1. Page path as Primary Key
2. Page Title
3. Page Content Snippet
4. PageRank
5. Created at
6. Updated at
7. IsHidden

##### **Link Table**
1. Page path PK
2. Link to page path FK
3. Created at
4. Updated at

##### **Term Table**
1. Page path PK
2. Term
3. Count
4. Created at
5. Updated at

##### **Image Table**
1. Page path PK
2. Image link
3. Fallback text
4. Created at
5. Updated at

### Redis
1. Used to enqueue to be indexed pages in a: "page path - page html" structure.
2. Used to track what pages link to other pages, have been crawled, and their Page rank in a:
"page path - list of connected pages - page rank" structure.
3. On 2. at start load all of the urls already in the PostgresDB, to prevent repetition.
If at memory limit for Redis, or as fallback make use of PostgresDB directly.

---

## User Web Interface
The User Web Interface could, in principle, work without the "Explore and Index Services",
as long as the PostgresDB holding the Indexed websites is also online.

### Query Engine
1. Written in Java Spring. Is in a Docker cluster managed by The Scaling Service.
2. Serves the Wengine's front page, and handles queries to the PostgresDB.
3. The front end is written in just HTML, CSS and JS.
4. Is a separate jar from The Scaling Service.
5. Available in a Docker image form. To be run in a Docker cluster, with a load balancer.
6. Hopefully will run it in the cloud for Demonstration purposes.

---

## Administration and Scaling
The idea is to require just the python cli, and docker. The rest should be in docker image form,
and be easily managed using either the cli or an administration website.

### Administration and Scaling Service
1. Written in Java Spring. Is controlled using http. Is separate from either User Web Interface or Explore and Index Services.
Can connect to many Docker daemons at once, and authenticates itself to them.
2. Makes sure Redis does not use too much memory
by scaling down Crawler and Indexer activity when above a certain threshold.
3. Automatically scales Crawler, Indexer, PageRank,
and Query Engine clusters, by talking to Docker, using docker-java library.
4. Listens for https commands in real time, and adjusts its activities based on them.
5. Can host a website that allows graphical administration (requiring auth of course).
The website is just HTML, CSS and JS.
6. Is in a separate jar and docker image from The Query Engine.
7. Can be used to start any and all parts of the project using https requests.
8. Available in a Docker image form, for ease of use.

### CLI Scaler Interface
1. Written in Python.
2. As the name suggests: is used to control the Scaling Service via Command Line Interface.
3. Connects to the Java Spring Scaling Service.
4. This is the recommended way to startup any and all parts of the project.
5. Should have a way to easily download all the needed
dependencies (most likely just docker to run all the images)
and docker images and running them. For ease of presentation.

---

## Additionally
1. Use automatic documentation libraries for all languages.
Preferring MarkDown document generation, but html is acceptable as an alternative.
2. Write Unit tests, where sensible. Same for Integration tests.
3. Use relevant language standards (eg. PEP8)
4. Use git and github. Write commits in the "Conventional Commits" standard.
