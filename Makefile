serve:
	docker build -t test . && docker run -it -p 8080:8080 -v $$(pwd):/app test 