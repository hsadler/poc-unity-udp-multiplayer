build:
	docker build -t poc-unity-udp-multiplayer:dev -f Dockerfile .

up:
	docker-compose -f docker-compose.yaml up

shell:
	docker exec -it poc-unity-udp-multiplayer /bin/bash