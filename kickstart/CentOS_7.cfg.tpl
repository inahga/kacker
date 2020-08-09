{% extends 'CentOS_8.cfg.tpl' %}

{% block authconfig %}
auth --enableshadow --passalgo=sha512
{% endblock %}

{% block packages %}
%packages
@^minimal
@core
perl
{% for package in custom_packages %}
{{ package }}
{% endblock %}
%end
{% endblock %}
