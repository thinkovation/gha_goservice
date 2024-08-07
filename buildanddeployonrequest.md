# Build and deploy on request
This action runs on request. It builds the app, and then copies it to a remote server and then restarts it.

## Setup
To set up for this action you need to use (or create) the following:

- A "deploy" user - which I usually call "deploy"
- A service control file for the application so it can be started/restarted remotely
- Sudo (no passwd) privs for your deployment user to invoke the start/restart commands
- A public/private keypair for the deployment user - the private key is used by github, the public key is put into the authorized keys for the deployment user

You will need to set up thr following secrets in the github repository:
- SSH_PRIVATE_KEY -  This is the private key you created for the deploy user
- REMOTE_USER: - The user name of the remote user (eg "deploy")
- REMOTE_HOST: - The host name / ip of the host you're deploying to
- REMOTE_PATH: - The remote path where you're deploying the app (This should match the path contained within the service definition file - see below)
- APP_NAME - the name of the app which will be used to name the binary, and the service file - in this example it is "goservice"

## Creating the service control file
The github action uses systemctl commands to start/restart the service - so you need to create a service file called {app_name}.service - So in our example it will be "goservice.service". You will need to replace every instance of "goservice" in the example below with your app name.

```
[Unit]
Description=GoService
After=network.target

[Service]
Type=simple
ExecStart=/home/deploy/goservice/goservice
WorkingDirectory=/home/deploy/goservice
Restart=on-failure
RestartSec=5s
User=deploy

[Install]
WantedBy=multi-user.target
```

This service definition file will automatically run the service after a reboot (specifically after the network is available). It will restart on failures, and it uses the deploy user (which is decent practice)

Once the file is created, it should be enabled - this command creates the necessary symlinks so that the service is started automatically during the boot process. Without enabling the service, it will need to be manually started after each reboot.

```
sudo systemctl enable goservice.service
```
For reference, you can disable it with the following command:
```
sudo systemctl disable goservice.service
```
This wont stop it - but it will prevent it from launching at the next restart



## Creating the user

In  this case I am calling the user "deploy". The user does not need a password but if you want to limit the ability of users without sudo privs to su to that user you may want to add a password.

```
sudo adduser deploy
```

This command will create a /home/deploy directory - Which is the directory we'll use to deploy the app. If you want to use another directory - for example opt/goservice then you'll need to create that directory and give the deploy user the relevant rwx access rights to it.

**Note - we're assuming that your server ssh config includes the "PasswordAuthentication no" directive - since we really don't want to allow password based access to our servers**

## Add SUDO permissions to user to allow them to start and stop your app.

You need to use the visudo command to open the editor for the sudoers file - This will open whichever editor has been configured as the default for modding the sudoers file.

Add the following lines (not forgetting to replace "goservice" with the service name you chose when you created the service control file)...
```
deploy ALL=(ALL) NOPASSWD: /bin/systemctl stop goservice.service
deploy ALL=(ALL) NOPASSWD: /bin/systemctl start goservice.service

```

**Explanation**
*deploy: The username to which you are granting permissions.
*ALL: Specifies that this rule applies to all hosts.
*(ALL): Specifies that the user can run the command as any user.
*NOPASSWD:: Specifies that no password is required.
*/bin/systemctl stop goservice.service: The exact command that can be run without a password.
*/bin/systemctl start goservice.service: The exact command that can be run without a password.

Ensure that the paths to systemctl and the service name are correct. The path /bin/systemctl is commonly used, but it might differ on some systems (it could be /usr/bin/systemctl on some distributions).


## Generate SSH Keypair for your user

Run ssh-keygen on either your local machine (probably best) or the server (but don't confuse it with a local key pair - this is the key pair we will use to connect to the server).

You will be asked for a name for the output files - I tend to use the name of the application, that way I can keep a track who  the keys are for. I also use the -C directive to put the app name into  the generated keys so you can tell from them what they relate to.

```
ssh-keygen -C goservice

```

Take the public key that is generated and put it into the authorized_keys file for the deploy user. 
In the next step you'll take the private key that is generated and add it as a repository secret called "SSH_PRIVATE_KEY".

## Set up secrets on your gitub repository

Go to your repository, select settings.

 Then on the side menu go to - Security - Secrets and Variables - Actions. Here you'll find your repository secrets and a button to add a new repository secret.

 You need to add the following (here the values I am using relate to my setup - modify them for yours)...

  SSH_PRIVATE_KEY -  take the text of the private key you generated in the previous step and paste it in here
- REMOTE_USER: - deploy
- REMOTE_HOST: - app.retrro.digital
- REMOTE_PATH: - /home/deploy
- APP_NAME - goservice

## Add the github workflow YAML to your repository

These files are stored in .github/workflows

So in my case I put the following into a file called .github/workflows/buildanddeployonrequest.yaml

```
name: Build and Deploy Go App

on:
  workflow_dispatch:

jobs:
  build_and_deploy:
    runs-on: ubuntu-latest

    env:
      SSH_PRIVATE_KEY: ${{ secrets.SSH_PRIVATE_KEY }}  # Add your private key to GitHub secrets
      REMOTE_USER: ${{ secrets.REMOTE_USER }}  # Add remote user to GitHub secrets
      REMOTE_HOST: ${{ secrets.REMOTE_HOST }}  # Add remote host to GitHub secrets
      REMOTE_PATH: ${{ secrets.REMOTE_PATH }}  # Add remote path to GitHub secrets
      APP_NAME: ${{ secrets.APP_NAME }}  # Add App name to GitHub secrets

    steps:
    - name: Checkout code
      uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.19'  # Replace with your desired Go version

    - name: Build
      run: |
        go build -o $APP_NAME ./cmd/$APP_NAME  # Adjust the command according to your project structure

    - name: Copy binary to remote server
      run: |
        mkdir -p ~/.ssh
        echo "$SSH_PRIVATE_KEY" > ~/.ssh/id_rsa
        chmod 600 ~/.ssh/id_rsa

        scp -o StrictHostKeyChecking=no ./$APP_NAME $REMOTE_USER@$REMOTE_HOST:$REMOTE_PATH/$APP_NAME

    - name: Restart service on remote server
      run: |
        ssh -o StrictHostKeyChecking=no $REMOTE_USER@$REMOTE_HOST << EOF
          sudo systemctl stop $APP_NAME.service  # Adjust this to your service name
          sudo cp $REMOTE_PATH/$APP_NAME /usr/local/bin/$APP_NAME  # Adjust this path to your binary location
          sudo systemctl start $APP_NAME.service  # Adjust this to your service name
EOF

```

### Key Points

#### Environment Variable `APP_NAME`:

- The `APP_NAME` environment variable is set to the name of your application ("myapp" in this example).
- This variable is used to dynamically name the executable file, specify the remote path, and restart the service.

#### Build Step:

- The `go build` command uses the `APP_NAME` variable to name the executable file.

#### Copy Binary Step:

- The `scp` command uses the `APP_NAME` variable to specify the destination path on the remote server.

#### Restart Service Step:

- The SSH commands use the `APP_NAME` variable to stop and start the service on the remote server and to copy the executable to the correct location.

### Using the Workflow

When you run this workflow, it will build the Go application, copy the resulting binary to the remote server, and restart the service using the specified `APP_NAME`. This makes the workflow more flexible and easier to maintain if the application name changes in the future.
