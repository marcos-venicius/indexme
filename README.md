# IndexMe

Pass a directory to the program then index all the sub directory files.

The first implementation will be using only TF-IDF technique.

Every "search" will be saved in a local sqlite file after do the indexing, preventing it to do the indexing everytime.

Then, the user could search upon a saved directory index and search for a file with a certain content.

The tool should keep a checksum for every file in the folder, this allows the user to call a flag `-update` and do a search
in all files of all subfolders and check if the hash matches, if yes, keep going if not re-index this file and only update them.

We also can cache the result of a search to improve the speed of the search, but everytime a file is updated we need to check if some cached search include this file, if yes, we remove this cache
to guarantee that the cached data is not outdated.

Default ignore folders like: `node_modules`, `.git`, ...
Default ignore files like: `binary files`, `image files`, `pdf files`, `data files`, `zip files`, ...

Allow the user to update this config by updating a configuration file (possible json).

## Examples

**Indexing:**

```bash
go run . -index /etc

# output something like

Indexing /etc/hosts...
Indexing /etc/resolv.conf...
Indexing /etc/foo...
Indexing /etc/foo/bar...
Indexing /etc/foo/bar/baz...

/etc indexed successfully
```

**Viewing indexed folders**

```bash
go run . -list

# output something like

/etc
  234 documents indexed
  20 ignored files (binary, images, any non readable file)
  last update at 2024-10-10 13:34 PM

/work/projects/todo-list
  234 documents indexed
  20 ignored files (binary, images, any non readable file)
  last update at 2024-10-10 13:34 PM
```

**Updating indexes**

```bash
go run . -update /etc

# output something like

/etc/resolv.conf already updated
/etc/hosts updated successfully

/etc was sucessfully updated
```

**Search for a term**

in the search, everything after the `-search` flag is a string to the query

```bash
go run . -search openssl config

# output something like

/etc/openssl/somefile.conf
/etc/other/file
/etc/foo/file
/etc/bar
/home/tests
```

or search in a specific directory:


```bash
go run . -dir /etc -search openssl config

# output something like

/etc/openssl/somefile.conf
/etc/other/file
/etc/foo/file
/etc/bar
```

**Removing folder**

```bash
go run . -remove /etc

# output something like

/etc folder removed sucessfully
```
