<img src="https://cdn-icons-png.flaticon.com/512/9227/9227748.png" alt="holocron" width="130"/>

# Holocron
Holocron is an object storage based leader election library.

Since Holocron stores everything in object storage that can be geo-replicated a node truly becomes virtual. In the future replacing a node does not require bootstrapping a new node. Simply pointing the new compute at the crashed nodes S3 `bucket+prefix` will allow it to resume immediately.

# Compiling
```
make
```

# Running the simulator
**_WARNING: THIS WILL CAUSE WRITES IN YOUR S3 BUCKET_**

Holocron doesn't build node binaries yet but at the moment a 2 node simulator is provided.
```
./build/simulator-arm64-darwin \
    -key=<REDACTED> \
    -secret=<REDACTED> \ 
    -region=<REDACTED> \
    -bucket=<REDACTED> \
    -prefix=leader
```
```
1:33PM INF node/node.go:70 > Starting node epoch=0 node=3f10d2708fce role=solo

1:33PM INF node/node.go:70 > Starting node epoch=0 node=d499b66c65f4 role=solo

1:33PM INF node/node.go:111 > Trying to acquire cluster leadership. epoch=0 node=3f10d2708fce role=solo

1:33PM INF node/node.go:111 > Trying to acquire cluster leadership. epoch=0 node=d499b66c65f4 role=solo
Blocking, press ctrl+c to continue...

1:33PM INF node/node.go:98 > Acquired cluster leadership. epoch=49 node=d499b66c65f4 role=leader

1:33PM INF node/node.go:89 > Failed to acquire cluster leadership, becoming a follower. epoch=49 node=3f10d2708fce role=follower

1:33PM INF node/node.go:111 > Trying to acquire cluster leadership. epoch=49 node=d499b66c65f4 role=leader

1:33PM INF node/node.go:111 > Trying to acquire cluster leadership. epoch=49 node=3f10d2708fce role=follower

1:33PM INF node/node.go:98 > Acquired cluster leadership. epoch=50 node=d499b66c65f4 role=leader

1:33PM INF node/node.go:89 > Failed to acquire cluster leadership, becoming a follower. epoch=50 node=3f10d2708fce role=follower
```

#### Ideas
- Implement virtual nodes with the ability to instantly resume from crashed nodes state.
- Each node can share an S3 `bucket+prefix` or each node can have it's own bucket or ACL's.

#### Current Features
- Stores all cluster metadata for each node in object storage.
- Uses object storage conditional writes to lock for leader election.
- Each epoch round creates a new lock file in an append-only approach.

#### TODO
- Conveniently store the current leaders view of the cluster topology in each epoch lock file so that we have a historical view from the leaders perspective.
- Stable leadership. At the moment every node competes every epoch on the leadership role. A more sophisticated algorithm that tries to keep a stable might be a good idea.
- Garbage collection of previous epoch lock files to prevent object storage costs from increasing and to prevent LIST performance from degrading to keep costs and performance constant.

# Status
Holocron is highly experimental and only an hour-ish of hacking. Basically all it does at the moment is compete on an epoch lock file. It doesn't implement lock expiration or anything that makes it a complete implementation yet.

# Contributing
If you're interested in helping out feel free! I'm curious where this will go.
