# How to host your server on AWS
## 1. Account
- Make sure to get a free tier
- Eligible for 6 months, you get $100 worth of credits
- Can also undergo a tutorial which provides you with another $100

## 2. Create EC2 Instance
1. Search `ec2` in the search bar (top left)
2. Click on `Launch Instance`
3. Assign any name
4. AMI (Amazon Machine Image): Choose `ubuntu` 
5. Instance Type: Choose `t3.micro`
6. Click on `Create new key pair`
    1. Enter key name
    2. Key pair type: `RSA`
    3. Private key file format: `.pem`
    4. Click on `Create key pair`
    5. This will download a private_key_name.pem file
7. Network Settings: Click on `Edit`
    1. There should be a default `Inbound Security Group Rules` of ssh:22. If not create one with following details
        - Type: `ssh`
        - Protocol: `TCP` (default)
        - Port: `22` (default)
        - Source type: `Anywhere`
        - Source: `0.0.0.0/0` (default)
    2. Click on `Add security group rule`
        - Type: `Custom TCP`
        - Protocol: `TCP` (default)
        - Port: `8080`
        - Source type: `Anywhere`
        - Source: `0.0.0.0/0` (default)
8. All done, Click on `Launch instance`. You can access instances in the `Instances` section of `EC2`

## 3. Connecting to EC2 Instance
1. Open a bash (git bash works) terminal in directory where your private key exists
2. run 
```bash
chmod 400 your-key.pem
ssh -i your-key.pem ubuntu@<your-public-ip>
```
3. Public IP can be found under `Public ipv4 address` column in `Instances` section of `EC2`

4. You would get an output like this, type yes and enter
```bash
The authenticity of host '34.226.212.235 (34.226.212.235)' can't be established.
ED25519 key fingerprint is SHA256:g/2qdauA8dA2vvKNUqP/kPTUydL8/p21eRh6NYTpDUw.
This key is not known by any other names.
Are you sure you want to continue connecting (yes/no/[fingerprint])?
```
5. Now you are inside the ec2 instance. There might be a slight lag between you typing and the text appearing in CLI from here on out.

## 4. Start your server
1. Initial setup
```bash
sudo apt update
sudo apt install docker.io -y
sudo systemctl start docker
```
2. Pull docker image and run
```bash
docker pull <your-dockerhub-username>/<image-name>:<version>
docker run -d -p 8080:8080 <your-dockerhub-username>/<image-name>:<version>
```

## 5. Submit Public IP and Port to Protohackers
1. You can also connect from your local machine using `telnet` command on powershell
```powershell
telnet <your-public-ip-address> <port>
```
2. This connects and then you can start typing