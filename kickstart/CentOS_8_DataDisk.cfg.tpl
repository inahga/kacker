{% extends "CentOS_8.cfg.tpl" %}

{% block diskconf %}
zerombr
clearpart --all --drives=sda,sdb
part /boot --fstype=ext4 --size=500 --ondisk=sda
part pv.01 --size=1024 --grow --ondisk=sda
part pv.02 --size=1024 --grow --ondisk=sdb
volgroup vg1 pv.01
volgroup vg2 pv.02
logvol swap --vgname=vg1 --size=2048 --name=lv_swap
logvol / --vgname=vg1 --fstype=xfs --size=2048 --grow --name=lv_root
logvol {{ diskpath|default('/data') }} --vgname=vg2 --fstype=xfs --size=2048 --grow --name=lv_data
{% endblock %}
