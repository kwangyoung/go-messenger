#!/bin/bash

function quit {
#	docker-compose stop
#	docker-compose rm -f
	docker stop Database Goauth
	docker rm -f Database Goauth
	exit $1
}

#docker-compose up -d

docker rm -f Database Goauth

docker run -d --name Database --expose 3306 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=rootpass -e MYSQL_DATABASE=messenger -e MYSQL_USER=messenger -e MYSQL_PASSWORD=messenger mysql
docker run -d --name Goauth -p 9000:9000 --link="Database:db" usman/go-auth:develop -l debug -t console run --db-host db -p 9000


# Make sure containers are ready for the test
sleep 20

if [ "$(uname -s)" = "Darwin" ] ; then
	service_ip=$(boot2docker ip)
else 
	service_container=$(docker ps -a | awk '{ print $1,$2 }' | grep go-auth | awk '{print $1 }')
	service_ip=$(docker inspect -f '{{ .NetworkSettings.IPAddress }}' ${service_container})
fi

echo "Using Service IP $service_ip"


first=$(curl -i -silent -X PUT -d userid=USERNAME -d password=PASSWORD ${service_ip}:9000/user | grep "HTTP/1.1")
second=$(curl -i -silent -X PUT -d userid=USERNAME -d password=PASSWORD ${service_ip}:9000/user | grep "HTTP/1.1")

status_first=$(echo "$first" | cut -f 2 -d ' ')
status_second=$(echo "$second" | cut -f 2 -d ' ')

if [[ "$status_first" -ne 200 ]]; then
	echo "Fail: Expecting 200 OK for first user register got $status_first"
	quit 1
else
	echo "Pass: Register User"
fi


if [[ "$status_second" -ne 409 ]]; then
	echo "Fail: Expecting 409 OK for second user register got $status_second"
	quit 1
else 
	echo "Pass: Register User Conflict"
fi

quit 0


