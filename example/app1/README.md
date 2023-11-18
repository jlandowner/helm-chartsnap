# Example of snapshot testing for local Helm chart

This chart repository is default helm repository created by `helm create app1` command.

And add [`test`](test) directory to place test patterns of the chart.

There is the example test values:

- [`test_hpa_enabled.yaml`](test/test_hpa_enabled.yaml)
- [`test_ingress_enabled.yaml`](test/test_ingress_enabled.yaml)

Do snapshot with the specific test values 📸

```sh
helm chartsnap -c . -f test/test_ingress_enabled.yaml
```

Or do snapshot for all test values 📸

```sh
helm chartsnap -c . -f test/ # specify directory for -f
```

Probably you will see the failure that does not match the snapshot with the above commands.

Then, update the snapshot with `-u` options.

```sh
helm chartsnap -c . -f test/ -u
```
