# MonzoUrlCrawler

This project starts with a base URL and crawls all href links found within the page, it will crawl all reachable pages within the base done

Support Schems are 
1. http
2. https


2 Main Concurrency Patterns I use
1. WorkerPool -> This is added in order to make sure we don't spwan 10000 of go routines if a page has 10000's of links
2. OrDone -> Clean way to check if ctx.Done has been called when dealing with results channel 


Flow 
1. Fetch Pages via fetcher
2. Parse the html via htmlparser, in order to extract <a href = "links" > 
3. Push unique links into worker in order to crawl pages within the same domain
4. Keep the loop going untill jobs count > 0
5. Once no more jobs in the channel, exit cleanly via closing channels
6. Result is a array of unique urls found within the monzo domain

Key Points
1. Resolved link fragements are being set to empty as they really don't serve a purpose for our crawl
2. resolved link path is normalised to "/" this is so that if any of the pages have a link back to home which doesn't contain "/" then we don't crawl it again


Total Pages Found : 42011

If more time given 
1. Create CI via github actions 
2. Create docker mock server and run against that in order to verify behaviour
3. Add rate limiting to calls made to the monzo server, if left unbounded cloudflare could consider this as an attack and block the IP, a wrapper service around the fether service will be sufficent. 