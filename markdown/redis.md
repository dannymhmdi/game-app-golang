# Redis short hand guide (CLI & ZSET)
# redis CLI
```bash
start ) redis-cli
  
1)get key *

desc:get all keys

2)set gameapp 1

desc:set gameapp key to 

3)get gameapp

desc:get gameapp value

4)keys *

desc:get all keys
```

# zset data-type
  1-Set: It stores unique elements (no duplicates).

  2-Sorted: Each element has a score associated with it, and the elements are automatically sorted by their scores.
  Think of it like a leaderboard where:

- Each player (element) has a unique name.

- Each player has a score (e.g., points, time, etc.).

- The players are automatically ranked from the highest to the lowest score.

---

# Key Features of ZSETs:

 - Unique Elements: Like a regular set, a ZSET cannot have duplicate elements.

 - Scores: Each element has a numeric score (e.g., a float or integer).

 - Automatic Sorting: Elements are always sorted by their scores, from lowest to highest.

 - Fast Operations: You can quickly add, remove, or update elements, and Redis will keep the set sorted.

---

# Common Use Cases for ZSETs:
1- Leaderboards: Rank users by their scores (e.g., in a game).

2- Priority Queues: Manage tasks with different priorities.

3- Time Series Data: Store timestamps as scores and events as elements.

4- Ranking Systems: Rank items by popularity, price, etc.

---

# Basic Commands for ZSETs:

Here are some common commands to work with ZSETs:
# 1. Add element

Use <b>ZADD<b> to add elements with their scores:
```bash
ZADD leaderboard 100 "Alice"
ZADD leaderboard 200 "Bob"
ZADD leaderboard 150 "Charlie"
```
- This adds three players to the leaderboard ZSET with their respective scores.

#  2. Get Elements by Rank:

Use <b>ZRANGE<b> to get elements in a specific rank range (sorted by scores).
```bash
ZRANGE leaderboard 0 -1
ZRANGE myset 0 -1 WITHSCORES
```
- This returns all elements in the ZSET, sorted from lowest to highest score:
- This returns all elements in the ZSET, sorted from lowest to highest score with their scores:
```bash
1) "Alice"
2) "Charlie"
3) "Bob"

- withscore:
1) "Alice"
2) "100"
3) "Charlie"
4) "150"
5) "Bob"
6) "200"
```



# 3. Get Elements by Score Range:

Use <b>ZRANGEBYSCORE<b> to get elements within a specific score range.
- 
```bash
 ZRANGEBYSCORE leaderboard 100 200
```
- This returns elements with scores between 100 and 200:
```bash
1) "Alice"
2) "Charlie"
3) "Bob"
```

# 4. Get Rank of an Element:
Use <b>ZRANK<b> to get the rank (position) of an element.
```bash
ZRANK leaderboard "Bob"
```
- This returns 2 because "Bob" is the third element (0-based index).

# 5. Get Score of an Element:
Use <b>ZSCORE<b> to get the score of a specific element.
```bash
ZSCORE leaderboard "Charlie"
```
- This returns 150, the score of "Charlie".

# 6. Remove Elements:
Use <b>ZREM<b> to remove an element from the ZSET.
```bash
ZREM leaderboard "Alice"
```
- This removes "Alice" from the ZSET.

# 7. Get the Size of the ZSET:
Use <b>ZCARD<b> to get the number of elements in the ZSET.
```bash
ZCARD leaderboard
```
- This returns the total number of elements in the ZSET.



