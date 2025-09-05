package client
/*

connect> ask>
name, [worker, consumer]

[worker]
connection running --> 
-[client_to_server] updating system info only if worker
-[server_to_client] recieve code
 
[consumer]
connection running --> 
- check compute capacity 
- summit code to server for running

*/

func Connect(serverIP string, port string) {

	// connect to websocket server.

	// host communication routes [consumer]
		// - send code for execution
		// - get compute that can be used

}

