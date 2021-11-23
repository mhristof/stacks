# stacks

Handle PR stack updates in git.

## Example

If you repo looks like this

<!-- markdownlint-disable MD013 -->
```shell
* 46c1740 2 seconds ago (Mike Christof) (mhristof@gmail.com)       (HEAD -> feat1.1)   feat1.1: commiting
|
* 4170581 2 seconds ago (Mike Christof) (mhristof@gmail.com)         feat1.1: commiting
|
* 4272275 2 seconds ago (Mike Christof) (mhristof@gmail.com)       (feat1)   feat1: commiting
|
| * 77c2bed 2 seconds ago (Mike Christof) (mhristof@gmail.com)       (main)   main: commiting
|/
|
* b6423f1 2 seconds ago (Mike Christof) (mhristof@gmail.com)         main: commiting
|
* 915c18c 2 seconds ago (Mike Christof) (mhristof@gmail.com)         main: commiting
|
* 6cd3ab5 3 seconds ago (Mike Christof) (mhristof@gmail.com)         initial import
```
<!-- markdownlint-enable MD013 -->

it means that your feat1 and feat1.1 branches are not based from the latest main.

To solve this, you could start rebasing by your self, or you could try the dry
run of `stacks`

```shell
$ stacks.darwin rebase -n
rebasing (branch: feat1.1 ) (dry: true )
command: git checkout feat1
command: git rebase --onto main main@{1}
command: git checkout feat1.1
command: git rebase --onto feat1 feat1@{1}
```

which will print out what commands are going to be executed and then pull the
trigger with

```shell
$ stacks.darwin rebase
rebasing (branch: feat1.1 ) (dry: false )
command: git checkout feat1
command: git rebase --onto main main@{1}
command: git checkout feat1.1
command: git rebase --onto feat1 feat1@{1}
```

which will produce this `git log`

<!-- markdownlint-disable MD013 -->
```shell
* 00191f1 2 minutes ago (Mike Christof) (mhristof@gmail.com)       (HEAD -> feat1.1)   feat1.1: commiting
|
* bf846f6 2 minutes ago (Mike Christof) (mhristof@gmail.com)         feat1.1: commiting
|
* 2ccc6ea 2 minutes ago (Mike Christof) (mhristof@gmail.com)       (feat1)   feat1: commiting
|
* 77c2bed 2 minutes ago (Mike Christof) (mhristof@gmail.com)       (main)   main: commiting
|
* b6423f1 2 minutes ago (Mike Christof) (mhristof@gmail.com)         main: commiting
|
* 915c18c 2 minutes ago (Mike Christof) (mhristof@gmail.com)         main: commiting
|
* 6cd3ab5 2 minutes ago (Mike Christof) (mhristof@gmail.com)         initial import
```
<!-- markdownlint-enable MD013 -->
