# Setting up the remote server

The server setup involves the following steps:
- Create a deployment user (I usually call this user "deploy")  that can write the binary to the directory it needs to be copied to
- Create a systemd service file to start and restart the service - give it a meaningful service name
- Give the deployment user sudo (no password) permission to invoke "service {servicename} start" and "service {servicename} restart"