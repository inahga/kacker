# Kacker
A wrapper for HashiCorp Packer that seeks to conquer the repetition in managing multiple Packer config files and Kickstart files. It's primarily targeted for using RHEL Kickstart files, but it should be extensible enough to be used on other provisioning scripts such as Windows' unattend.xml or Debian preseed files.

## Customization
The key file for running Kacker is the Customization file. It resolves templated kickstart files and packer files, and runs Packer against them.

To run a customization:
```
kacker customization.yaml
```

To supply flags to Packer, such as a secret variable file:
```
kacker -packer-flags='-var-file=secrets.json' customization.yaml
```

See `examples/` for examples of customization files and templated kickstart files.

## Kickstart Templating
The `kickstart` section of the customization provides a hierarchy of templated kickstart files, and supplies variables to them.

```
kickstart:
  from:
    - kickstart_parent.cfg.tmpl
    - kickstart_child.cfg.tmpl
  variables:
    - name: var1
      value: var1
    - name: var2
      values: [array, of, values]
```

Variables can be of type `value`, `values`, `url`, `urls`, `fragment`, and `fragments`. The array variable types map to arrays in the template file. Non-array variable types map to single variables.

`url` and `urls` issue a GET request to the supplied URL and store the results in a variable or array, respectively. This is useful for grabbing plain text documents that exist over the internet. No assumptions about whitespace are made. No concatenation occurs.

`fragment` and `fragments` pull files from the file system into a variable or array, respectively. This is good for supplying several scripts to the kickstart file.

## Packer Config Files
Packer already uses the Go template engine for its own means. Rather than muddy the waters by running our own template engine, overriding of values is accomplished through YAML merging.

## Todo
- Multiple customizations in a single file with parallel execution.
- Enhanced image lifecycle management.
    - Packer will delete images before creating a new one, at least for vSphere. Could interrupt other builds that rely on an image.
- Add other YAML merge/override strategies.
- Enhance docs.
- Get unit test coverage >70%
