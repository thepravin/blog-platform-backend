
- composite indexing

- view architect
    - don't count when user press f5
    - if user comback after 24h cout it as a fresh view
    - every time updating the view count when fetching the get post by id it destory the primary database with the replica db.
    - shouldn't count own post view count.
    - from UI avoid 
        - debouncing (acdiently open the page)
        - press f5 multiple times
        - refresh page
    
    - soln: 
        - when user load page wait for 5 sec to check user dosen't accidentally click or bounce.
        - react send quite post request ot /posts/:id/view in background
        - store the history in post_view table and check bases on user ip or agent in the last 24 hours.
        - setup the cron job every 10 minutes to count the rows and update the main posts table in one big batch.


