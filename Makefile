up:
	docker-compose down --remove-orphans
	docker-compose up --watch --build
front: 
	docker-compose up --watch --build frontend