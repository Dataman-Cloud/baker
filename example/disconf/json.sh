#!/bin/bash

curl "http://192.168.1.21:8000/authorize" \
-H "Content-Type:application/json" \
--data @<(cat <<EOF
{
      "username":"admin",
      "password":"badmin",
}		
EOF
)



