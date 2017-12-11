# vstore_tester

## Test

```
gcloud beta emulators datastore start --no-store-on-disk

$(gcloud beta emulators datastore env-init)

./test.sh
```