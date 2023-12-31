# Hermetic

Sends data to digital storage

## Usage

### Send

```shell
hermetic send \
    --kafka-endpoints=<list-of-kafka-endpoints> \
    --transfer-topic <topic-name> \
    --stage-artifacts-root <stage-artifacts-root>
```

### Verify

```shell
hermetic verify \
    --kafka-endpoints=<list-of-kafka-endpoints> \
    --reject-topic <topic-name>
```

### Acquisition upload

```shell
hermetic acquistion \
    --kafka-endpoints=<list-of-kafka-endpoints> \
    --transfer-topic <topic-name> \
    --acquisition-root=</path/to/root/of/acquisition>
```
