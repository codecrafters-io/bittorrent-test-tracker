run_tracker:
	docker build -t test . && docker run -it -p 8080:8080 test /bin/tracker.sh

run_seeder:
	docker build -t test . && docker run -it -p 8080:8080 test /bin/seeder.sh

test_production:
	curl -G \
		--url-query "+info_hash=I%ED%8BH%C12%97LZ0%FC_%7Bn%1F%E7w7%A4%DF" \
		--url-query "+peer_id=I%ED%8BH%C12%97LZ0%FC_%7Bn%1F%E7w7%A4%DF" \
		--url-query "+left=0" \
		--url-query "+downloaded=0" \
		--url-query "+uploaded=0" \
		--url-query "+port=1234" \
		http://bittorrent-test-tracker.codecrafters.io/announce