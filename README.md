# go-cbox-scraper
Handles scraping from cbox.ws

Note: this package only handles "classic" cbox pages

Note 2: the entire scraper state is saved in the cache file
This means - debug level, box id, webhost id, etc.
If you need to change any of those values, simply load the cache, change the values you need, and save.
