#!/usr/bin/python3
#
import yaml, json, sys
print(json.dumps(yaml.load(sys.stdin.read(), Loader=yaml.FullLoader), indent=2))
