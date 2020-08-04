# ********* server gRCP-gw-REST *********
# server fires separate go routines for gRCP & REST gateway
# start server in separate terminal:
# change directory to user-service
# *** run server by :: go run .

# ********* client gRCP
# to run client exaqmples in go :: go run cmd/cli/*
# it will fire 	cli.ExampleRunClientGRCP() & cli.ExampleRunClientREST()

# ********* client REST by curl
# ********* simple test
# curl -X GET "http://localhost:8080/users"
# curl -X GET "http://localhost:8080/users?active=true"
# curl -X GET "http://localhost:8080/users?active=false"
# curl -X POST -d '{"active":true,"name":"Name-1 Last-T"}' "http://localhost:8080/users?active=true"
# curl -X POST -d '{"active":false,"name":"Name-2 Last-F"}' "http://localhost:8080/users?active=true"
# curl -X GET "http://localhost:8080/users"
# curl -X GET "http://localhost:8080/users?active=true"
# curl -X GET "http://localhost:8080/users?active=false"
# curl -X GET "http://localhost:8080/users?active="
# curl -X GET "http://localhost:8080/users?active=1"
# curl -X GET "http://localhost:8080/users?active=0"
# curl -X DELETE "http://localhost:8080/users/1001"
# curl -X GET "http://localhost:8080/users"
# curl -X GET "http://localhost:8080/users?active=true"
# curl -X GET "http://localhost:8080/users?active=false"
# curl -X DELETE "http://localhost:8080/users/1002"
# curl -X GET "http://localhost:8080/users"
# curl -X GET "http://localhost:8080/users?active=true"
# curl -X GET "http://localhost:8080/users?active=false"
# ---------------------------------
# 
# to clear storage db use truncate dev hack: curl -X GET "http://localhost:8080/users?active=dev-db-truncate"
# 
# 
# ********* gRCP-gw-REST client playground *********
# 
NODE="http://localhost:8080"
# NODE="http://192.168.1.53:8080"
# 
HEAD=""
# 
# for secure TLS (use in pair with secure HEAD below)
# NODE="https://localhost:8080"
#
# for secure TLS and auth
# HEAD=" --cacert "cert/server.crt" -H "sys:ss" -H "org:oo" -H "login:ll" -H "password:pp" "
# for secure TLS (no cert verification) and auth
# HEAD=" -k -H "sys:ss" -H "org:oo" -H "login:ll" -H "password:pp" "
# 
START=$(date +"%s%N")
clear
echo "------START------ ${START} >>> ${NODE} >>> "
echo && echo "--clr--" && curl ${HEAD} -X GET "${NODE}/users?active=dev-db-truncate"
echo && echo && echo "--ADD--&--LIST--"
echo && echo "-add-" && curl ${HEAD} -X POST -d '{"active":true}' "${NODE}/users"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-add-" && curl ${HEAD} -X POST -d '{"active":true,"name":"Name-1 Last-T"}' "${NODE}/users"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-add-" && curl ${HEAD} -X POST -d '{"active":false,"name":"Name-2 Last-F"}' "${NODE}/users"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-add-" && curl ${HEAD} -X POST -d '{"name":"Name-3 Last-F"}' "${NODE}/users"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-add-" && curl ${HEAD} -X POST -d '{"active":true,"name":"Name-4 Last-T"}' "${NODE}/users"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-add-" && curl ${HEAD} -X POST -d '{"active":false,"name":"Name-5 Last-F"}' "${NODE}/users"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-add-" && curl ${HEAD} -X POST -d '{"active":true,"name":"Name-6 Last-T"}' "${NODE}/users"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo && echo "--DELETE--&--LIST--"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-del-" && curl ${HEAD} -X DELETE "${NODE}/users/1001"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-del-" && curl ${HEAD} -X DELETE "${NODE}/users/1001"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-del-" && curl ${HEAD} -X DELETE "${NODE}/users/1002"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-del-" && curl ${HEAD} -X DELETE "${NODE}/users/1002"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-del-" && curl ${HEAD} -X DELETE "${NODE}/users/1003"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-del-" && curl ${HEAD} -X DELETE "${NODE}/users/1004"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo "-del-" && curl ${HEAD} -X DELETE "${NODE}/users/1005"
echo && echo "-lst-" && curl ${HEAD} -X GET "${NODE}/users"
echo && echo && echo "--DB--"
echo && echo "--clr--" && curl ${HEAD} -X GET "${NODE}/users?active=dev-db-truncate"
echo && echo "--mck--" && curl ${HEAD} -X GET "${NODE}/users?active=dev-db-mockup"
echo && echo "--shr--" && curl ${HEAD} -X GET "${NODE}/users?active=dev-db-shrink"
echo && echo "--clr--" && curl ${HEAD} -X GET "${NODE}/users?active=dev-db-truncate"
echo && echo "--mck--" && curl ${HEAD} -X GET "${NODE}/users?active=dev-db-mockup-100"
echo && echo "--clr--" && curl ${HEAD} -X GET "${NODE}/users?active=dev-db-truncate"
# 
END=$(date +"%s%N")
TNS=$(((${END}-${START})/1))
TUS=$(((${END}-${START})/1000))
TMS=$(((${END}-${START})/1000000))
TS=$(((${END}-${START})/1000000000))
# 
echo && echo && echo "------END------ ${END} >>> ${NODE} >>> (${TS} s, ${TMS} ms, ${TUS} us, ${TNS} ns)"
# 
# ********* playground end *********
# 
# -------- notes
# 
# REST GW insecure (no TLS handshake, no user/auth interceptor)
# memory
# ------END------ 1596370176324721217 >>> localhost >>> (0 s, 218 ms, 218584 us, 218584633 ns)
# ------END------ 1596370120923700532 >>> 192.168.1.53 >>> (0 s, 263 ms, 263097 us, 263097088 ns)
# persist
# ------END------ 1596370280772783537 >>> localhost >>> (2 s, 2516 ms, 2516014 us, 2516014986 ns)
# ------END------ 1596370405797706301 >>> 192.168.1.53 >>> (2 s, 2645 ms, 2645307 us, 2645307007 ns)
# 
# REST GW secure (with TLS handshake, with user/auth interceptor per call)
# memory
# ------END------ 1596372487554309238 >>> localhost >>> (0 s, 225 ms, 225900 us, 225900406 ns)
# persist
# ------END------ 1596372600525864317 >>> localhost >>> (2 s, 2514 ms, 2514646 us, 2514646711 ns)
# 
# 
# {"error":"strconv.ParseBool: parsing \"10\": invalid syntax","code":3,"message":"strconv.ParseBool: parsing \"10\": invalid syntax"} 
# {"error":"strconv.ParseBool: parsing \"-1\": invalid syntax","code":3,"message":"strconv.ParseBool: parsing \"-1\": invalid syntax"}
# {"error":"strconv.ParseBool: parsing \"yes\": invalid syntax","code":3,"message":"strconv.ParseBool: parsing \"yes\": invalid syntax"} 
# {"error":"strconv.ParseBool: parsing \"no\": invalid syntax","code":3,"message":"strconv.ParseBool: parsing \"no\": invalid syntax"}
# {"error":"strconv.ParseBool: parsing \"b\": invalid syntax","code":3,"message":"strconv.ParseBool: parsing \"b\": invalid syntax"}
# {"error":"strconv.ParseBool: parsing \"b\": invalid syntax","code":3,"message":"strconv.ParseBool: parsing \"b\": invalid syntax"}
# on secure and no credentials {"error":"unknown user ","code":2,"message":"unknown user "}
# 
#  TLS cert error, TODO new cert with ip
# ------START >>> 192.168.1.53 >>> 
# {"error":"connection error: desc = \"transport: authentication handshake failed: x509: cannot validate certificate for 192.168.1.53 because it doesn't contain any IP SANs\"","code":14,"message":"connection error: desc = \"transport: authentication handshake failed: x509: cannot validate certificate for 192.168.1.53 because it doesn't contain any IP SANs\""}
# 
# 