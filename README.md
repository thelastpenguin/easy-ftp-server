# Easy FTP Server by Gareth George
Easy FTP Server is designed to be an extremely light weight and easy to configure FTP Server that simply works right out of the box for most simple use cases.

# Configuration
create a file ~/.easyftp in the home directory of the user running the daemon with the following structure
```
{
    "/var/www/": [
        {
            username: "..."
            password: "..."
            read: true,
            write: true
        }
    ],
    "directory": [
        list of users
        ...
    ]
}
```
