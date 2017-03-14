# Easy FTP Server by Gareth George
Easy FTP Server is designed to be an extremely light weight and easy to configure FTP Server that simply works right out of the box for most simple use cases.

Note that there are obvious security implications since the FTP Server runs as whatever user is hosting it so all files it will creates will, by default,
also be owned by this user.

This is not a problem however for the usecases like a website with multiple subdomains and different user accounts managing each subdomain

# Configuration
create a file ~/.easyftp in the home directory of the user running the daemon with the following structure
```
{
    "Host": "127.0.0.1",
    "Port": 8021,
    "Users": [
        {
            "Username": "webadmin",
            "Password": "1234",
            "FsRoot": "/var/www/"
        },
		...
    ]
}
```
