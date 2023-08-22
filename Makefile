run_tracker:
	docker build -t test . && docker run -it -p 8080:8080 test /bin/tracker.sh

test_production:
	# curl -G --url-query "+info_hash=I%ED%8BH%C12%97LZ0%FC_%7Bn%1F%E7w7%A4%DF" http://bittorrent-test-tracker.codecrafters.io/scrape
	curl -G --url-query "+info_hash=I%ED%8BH%C12%97LZ0%FC_%7Bn%1F%E7w7%A4%DF" https://bittorrent-test-tracker-411072e35b04.herokuapp.com/scrape