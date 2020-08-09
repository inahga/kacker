# Kacker
This project seeks to reduce the ridiculous amount of repetition when using Packer to build images. Right now, it's targeted towards resolving RHEL kickstart files, but should be extensible enough to handle other config types (even Windows unattend.xml).


### Old documentation 
Using this repo involves the following command:

```
make TARGET=<packerfile> KICKSTART=<kickstart file> OVERRIDE=<override args> <command>
```
where `<packerfile>` is the packerfile in YAML format. It will automatically be converted to JSON.

where `<kickstart>` is the name of the kickstart file, relative to the `http` directory. This is not always required, it depends on the target packerfile. Windows will use the same directive, despite it being called an unattend/answer file.

where `<override args>` is a list of arguments to pass directly to Packer. This entire directive can be omitted, if not needed.

where `<command>` is one of the following:
 - "echo": print the contents of the packerfile to STDOUT in JSON format.
 - "validate": corresponds to `packer validate`
 - "build": corresponds to `packer build`
 - "debug": corresponds to `packer build -debug`

Example commands:
 - `make TARGET=CentOS_8.1.yml KICKSTART=CentOS_8.1_ks.cfg OVERRIDE='-var "CPUs=4" -var "RAM=4096"' build`
 - `make TARGET=CentOS_8.1.yml KICKSTART=CentOS_8.1_ks.cfg validate`

You must have `packer` in your `$PATH` and Python 3 installed.

You should rename `secrets.json.template` to `secrets.json` and populate it with info.

### HTTP Server
Some Linux distros have removed floppy disk support (e.g. CentOS/RHEL 8). The workaround is to use the Packer local HTTP server. This will require you to open ports on your local machine's firewall (applicable for RHEL distros--Ubuntu firewall is off by default).

You can turn the firewall off with `sudo systemctl stop firewalld`, or open a port range with `sudo firewall-cmd --zone=public --add-port=8000-9000/tcp`.


### Kickstart Files
To validate kickstart files, install `pykickstart` and use the `ksvalidator` command.
