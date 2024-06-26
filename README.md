

Based on the requirement, here are the entities I should have on the database at least:

	a. User,
	b. Subscription (including available packages)
	c. Activities (like & pass)

Below, I choose to have 3 database engines: MongoDB, PostgreSQL, and Redis.
	
MongoDB
I will use MongoDB to store activities. Since ‘activities’ is the main feature of a dating
app, we must anticipate a high read/write and a fast-growing data size on its storage
– NoSQL was made for that.
	
PostgreSQL
This SQL engine, with its attribute of ACID, is suitable for storing user and subscription data.


Redis

Caching using Redis avoids the frequent pulling of data to the main database.

	1. storing users' subscriber status,
	2. storing profile details for feeds
	3. plus, serving a locking mechanism (redsync) for subscriptions, preventing race conditions while updating user activities
	
Redsync

In case you deploy multiple machines/containers per service, a distributed locking
mechanism is required.
