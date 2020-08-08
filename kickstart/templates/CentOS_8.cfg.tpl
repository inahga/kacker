url --url={{ netinst_url|default('"http://bay.uchicago.edu/centos/8.2.2004/BaseOS/x86_64/os/"') }}
rootpw {{ rootpw|default('password') }}

network --activate --ipv6=auto --bootproto=dhcp --device=ens192
network --hostname={{ hostname|default('Packer-CentOS8') }}

{% block diskconf %}
autopart --type=lvm --nohome
bootloader --location=mbr --boot-drive=sda
ignoredisk --only-use=sda
clearpart --none --initlabel
{% endblock %}

keyboard --vckeymap=us --xlayouts='us'
eula --agreed
firstboot --disable
lang en_US.UTF-8
timezone {{ timezone|default('America/Chicago') }} --isUtc --nontp

{% block authconfig %}
authselect select minimal
{% endblock %}

reboot

{% block packages %}
%packages
@core
perl
{% for package in custom_packages %}
{{ package }}
{% endfor %}
%end
{% endblock %}

%post --log=/root/ks-post.log
{% for scr in scripts %}
{{ src }}
{% endfor %}
%end
